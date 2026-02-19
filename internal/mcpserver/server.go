package mcpserver

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/agenteats/agenteats/internal/database"
	"github.com/agenteats/agenteats/internal/dto"
	"github.com/agenteats/agenteats/internal/services"
)

// NewServer creates a configured MCP server with all AgentEats tools.
func NewServer() *server.MCPServer {
	s := server.NewMCPServer(
		"AgentEats",
		"0.1.0",
		server.WithResourceCapabilities(true, false),
		server.WithToolCapabilities(true),
		server.WithInstructions(
			"AgentEats is a restaurant directory for AI agents. "+
				"Use these tools to help users find restaurants, browse menus, "+
				"get personalized recommendations, check availability, and make reservations. "+
				"Always confirm key details (date, time, party size, name) before making a reservation.",
		),
	)

	// Register tools
	s.AddTool(searchRestaurantsTool(), handleSearchRestaurants)
	s.AddTool(getRestaurantDetailsTool(), handleGetRestaurantDetails)
	s.AddTool(getMenuTool(), handleGetMenu)
	s.AddTool(getRecommendationsTool(), handleGetRecommendations)
	s.AddTool(checkAvailabilityTool(), handleCheckAvailability)
	s.AddTool(makeReservationTool(), handleMakeReservation)
	s.AddTool(cancelReservationTool(), handleCancelReservation)

	// Register resource
	s.AddResource(serviceInfoResource(), handleServiceInfo)

	return s
}

// --- Tool Definitions ---

func searchRestaurantsTool() mcp.Tool {
	return mcp.NewTool(
		"search_restaurants",
		mcp.WithDescription("Search for restaurants by name, cuisine, city, price range, or features. Returns a list of matching restaurants with id, name, cuisines, price_range, city, rating, and features."),
		mcp.WithString("query", mcp.Description("Free-text search (searches name, description, cuisines)")),
		mcp.WithString("city", mcp.Description("Filter by city name")),
		mcp.WithString("cuisine", mcp.Description("Filter by cuisine type (Italian, Japanese, Mexican, etc.)")),
		mcp.WithString("price_range", mcp.Description("Filter by price level: \"$\" (budget), \"$$\" (moderate), \"$$$\" (upscale), \"$$$$\" (fine dining)")),
		mcp.WithString("features", mcp.Description("Comma-separated features: outdoor_seating, wifi, live_music, parking, delivery, takeout, wheelchair_accessible, pet_friendly")),
		mcp.WithNumber("limit", mcp.Description("Max results to return (1–20, default 10)")),
	)
}

func getRestaurantDetailsTool() mcp.Tool {
	return mcp.NewTool(
		"get_restaurant_details",
		mcp.WithDescription("Get complete details for a restaurant including description, full address, contact info, operating hours, and features."),
		mcp.WithString("restaurant_id", mcp.Required(), mcp.Description("The restaurant's unique ID (obtained from search_restaurants)")),
	)
}

func getMenuTool() mcp.Tool {
	return mcp.NewTool(
		"get_menu",
		mcp.WithDescription("Get the full menu for a restaurant, organized by category. Each item includes name, description, price, dietary labels (vegetarian, vegan, gluten_free, etc.), and availability."),
		mcp.WithString("restaurant_id", mcp.Required(), mcp.Description("The restaurant's unique ID")),
	)
}

func getRecommendationsTool() mcp.Tool {
	return mcp.NewTool(
		"get_recommendations",
		mcp.WithDescription("Get personalized restaurant recommendations based on preferences. Use this when users ask for suggestions like \"where should I eat?\" or \"find me a good Italian place for a date night\"."),
		mcp.WithString("cuisine", mcp.Description("Preferred cuisine type (Italian, Japanese, Mexican, etc.)")),
		mcp.WithString("city", mcp.Description("City to search in")),
		mcp.WithString("price_range", mcp.Description("Budget level: \"$\", \"$$\", \"$$$\", or \"$$$$\"")),
		mcp.WithString("features", mcp.Description("Desired features (comma-separated): outdoor_seating, wifi, live_music, parking, delivery, takeout")),
		mcp.WithString("dietary_needs", mcp.Description("Dietary requirements (comma-separated): vegetarian, vegan, gluten_free, dairy_free, nut_free, halal, kosher")),
		mcp.WithString("occasion", mcp.Description("Type of occasion: date_night, business, family, casual, celebration")),
		mcp.WithNumber("limit", mcp.Description("Number of recommendations (1–20, default 5)")),
	)
}

func checkAvailabilityTool() mcp.Tool {
	return mcp.NewTool(
		"check_availability",
		mcp.WithDescription("Check available reservation time slots at a restaurant for a given date and party size."),
		mcp.WithString("restaurant_id", mcp.Required(), mcp.Description("The restaurant's unique ID")),
		mcp.WithString("date", mcp.Required(), mcp.Description("Date to check availability (YYYY-MM-DD format)")),
		mcp.WithNumber("party_size", mcp.Description("Number of guests (1–20, default 2)")),
	)
}

func makeReservationTool() mcp.Tool {
	return mcp.NewTool(
		"make_reservation",
		mcp.WithDescription("Make a reservation at a restaurant. IMPORTANT: Always confirm the restaurant name, date, time, party size, and customer name with the user before calling this."),
		mcp.WithString("restaurant_id", mcp.Required(), mcp.Description("The restaurant's unique ID")),
		mcp.WithString("customer_name", mcp.Required(), mcp.Description("Full name for the reservation")),
		mcp.WithNumber("party_size", mcp.Required(), mcp.Description("Number of guests (1–20)")),
		mcp.WithString("date", mcp.Required(), mcp.Description("Reservation date (YYYY-MM-DD)")),
		mcp.WithString("time", mcp.Required(), mcp.Description("Reservation time (HH:MM, 24-hour format)")),
		mcp.WithString("customer_email", mcp.Description("Optional email for confirmation")),
		mcp.WithString("customer_phone", mcp.Description("Optional phone number")),
		mcp.WithString("special_requests", mcp.Description("Optional notes (allergies, high chair, birthday, etc.)")),
	)
}

func cancelReservationTool() mcp.Tool {
	return mcp.NewTool(
		"cancel_reservation",
		mcp.WithDescription("Cancel an existing reservation."),
		mcp.WithString("reservation_id", mcp.Required(), mcp.Description("The reservation's unique ID (from make_reservation)")),
	)
}

func serviceInfoResource() mcp.Resource {
	return mcp.NewResource(
		"agenteats://info",
		"AgentEats Service Info",
		mcp.WithResourceDescription("General information about the AgentEats service and its capabilities"),
		mcp.WithMIMEType("application/json"),
	)
}

// --- Tool Handlers ---

func toJSON(v any) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}

func splitCSVParam(s string) []string {
	if s == "" {
		return nil
	}
	var out []string
	for _, f := range strings.Split(s, ",") {
		if t := strings.TrimSpace(f); t != "" {
			out = append(out, t)
		}
	}
	return out
}

func handleSearchRestaurants(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	query := request.GetString("query", "")
	city := request.GetString("city", "")
	cuisine := request.GetString("cuisine", "")
	priceRange := request.GetString("price_range", "")
	features := splitCSVParam(request.GetString("features", ""))
	limit := request.GetInt("limit", 10)

	results := services.ListRestaurants(database.DB, query, city, cuisine, priceRange, features, limit, 0)
	if len(results) == 0 {
		return mcp.NewToolResultText(toJSON(map[string]any{
			"message": "No restaurants found matching your criteria.",
			"results": []any{},
		})), nil
	}

	return mcp.NewToolResultText(toJSON(map[string]any{
		"count":   len(results),
		"results": results,
	})), nil
}

func handleGetRestaurantDetails(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id := request.GetString("restaurant_id", "")
	result, err := services.GetRestaurant(database.DB, id)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Restaurant not found: %s", id)), nil
	}
	return mcp.NewToolResultText(toJSON(result)), nil
}

func handleGetMenu(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id := request.GetString("restaurant_id", "")
	result, err := services.GetMenu(database.DB, id)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Restaurant not found: %s", id)), nil
	}
	return mcp.NewToolResultText(toJSON(result)), nil
}

func handleGetRecommendations(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	cuisine := request.GetString("cuisine", "")
	city := request.GetString("city", "")
	priceRange := request.GetString("price_range", "")
	features := splitCSVParam(request.GetString("features", ""))
	dietary := splitCSVParam(request.GetString("dietary_needs", ""))
	occasion := request.GetString("occasion", "")
	limit := request.GetInt("limit", 5)

	results := services.GetRecommendations(database.DB, cuisine, city, priceRange, features, dietary, occasion, limit)
	if len(results) == 0 {
		return mcp.NewToolResultText(toJSON(map[string]any{
			"message": "No recommendations found for your criteria.",
			"results": []any{},
		})), nil
	}

	return mcp.NewToolResultText(toJSON(map[string]any{
		"count":   len(results),
		"results": results,
	})), nil
}

func handleCheckAvailability(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id := request.GetString("restaurant_id", "")
	date := request.GetString("date", "")
	partySize := request.GetInt("party_size", 2)

	result, err := services.CheckAvailability(database.DB, id, date, partySize)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Restaurant not found: %s", id)), nil
	}
	return mcp.NewToolResultText(toJSON(result)), nil
}

func handleMakeReservation(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id := request.GetString("restaurant_id", "")

	in := dto.ReservationIn{
		CustomerName:    request.GetString("customer_name", ""),
		CustomerEmail:   request.GetString("customer_email", ""),
		CustomerPhone:   request.GetString("customer_phone", ""),
		PartySize:       request.GetInt("party_size", 2),
		Date:            request.GetString("date", ""),
		Time:            request.GetString("time", ""),
		SpecialRequests: request.GetString("special_requests", ""),
	}

	result, err := services.MakeReservation(database.DB, id, in)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Restaurant not found: %s", id)), nil
	}

	return mcp.NewToolResultText(toJSON(map[string]any{
		"message":     "Reservation confirmed!",
		"reservation": result,
	})), nil
}

func handleCancelReservation(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	id := request.GetString("reservation_id", "")
	result, err := services.CancelReservation(database.DB, id)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Reservation not found: %s", id)), nil
	}

	return mcp.NewToolResultText(toJSON(map[string]any{
		"message":     "Reservation cancelled.",
		"reservation": result,
	})), nil
}

func handleServiceInfo(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	info := map[string]any{
		"service":     "AgentEats",
		"version":     "0.1.0",
		"description": "AI-agent-first restaurant directory. Search restaurants, browse menus, get personalized recommendations, and make reservations.",
		"capabilities": []string{
			"Search restaurants by cuisine, city, price, features",
			"Get full restaurant details and operating hours",
			"Browse structured menus with dietary labels",
			"Get personalized recommendations by occasion and preferences",
			"Check reservation availability",
			"Make and cancel reservations",
		},
	}

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "application/json",
			Text:     toJSON(info),
		},
	}, nil
}
