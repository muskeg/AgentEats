package main

import (
	"fmt"
	"log"

	"github.com/agenteats/agenteats/internal/config"
	"github.com/agenteats/agenteats/internal/database"
	"github.com/agenteats/agenteats/internal/models"
)

func ptr(f float64) *float64 { return &f }
func intPtr(i int) *int      { return &i }

type menuEntry struct {
	Category    string
	Name        string
	Description string
	Price       float64
	Labels      string
	Popular     bool
	Calories    *int
}

type restaurantSeed struct {
	Info         models.Restaurant
	WeekdayOpen  string
	WeekdayClose string
	WeekendOpen  string
	WeekendClose string
	ClosedDays   []string
	Menu         []menuEntry
}

var days = []string{"monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"}

func buildHours(restaurantID, wdOpen, wdClose, weOpen, weClose string, closed []string) []models.OperatingHours {
	if wdOpen == "" {
		wdOpen = "11:00"
	}
	if wdClose == "" {
		wdClose = "22:00"
	}
	if weOpen == "" {
		weOpen = "10:00"
	}
	if weClose == "" {
		weClose = "23:00"
	}

	closedSet := make(map[string]bool)
	for _, d := range closed {
		closedSet[d] = true
	}

	hours := make([]models.OperatingHours, 0, 7)
	for _, day := range days {
		isWeekend := day == "saturday" || day == "sunday"
		o, c := wdOpen, wdClose
		if isWeekend {
			o, c = weOpen, weClose
		}
		hours = append(hours, models.OperatingHours{
			RestaurantID: restaurantID,
			Day:          day,
			OpenTime:     o,
			CloseTime:    c,
			IsClosed:     closedSet[day],
		})
	}
	return hours
}

func getSeedData() []restaurantSeed {
	return []restaurantSeed{
		// 1. Bella Notte
		{
			Info: models.Restaurant{
				Name:        "Bella Notte",
				Description: "Authentic Italian trattoria with handmade pasta, wood-fired pizza, and an extensive wine list. Intimate candlelit atmosphere perfect for date night.",
				Cuisines:    "Italian,Mediterranean",
				PriceRange:  models.PriceUpscale,
				Address:     "142 Thompson St",
				City:        "New York",
				State:       "NY",
				ZipCode:     "10012",
				Country:     "US",
				Latitude:    ptr(40.7270),
				Longitude:   ptr(-73.9990),
				Phone:       "+1-212-555-0142",
				Email:       "reservations@bellanotte.example.com",
				Website:     "https://bellanotte.example.com",
				Features:    "outdoor_seating,wifi,live_music,wheelchair_accessible",
				TotalSeats:  80,
				Rating:      ptr(4.7),
				ReviewCount: 342,
			},
			WeekdayOpen: "17:00", WeekdayClose: "23:00", WeekendOpen: "12:00", WeekendClose: "00:00",
			Menu: []menuEntry{
				{"Appetizer", "Bruschetta Trio", "Tomato basil, mushroom truffle, and nduja spread on grilled sourdough", 16.0, "vegetarian", true, nil},
				{"Appetizer", "Burrata Caprese", "Fresh burrata with heirloom tomatoes, basil oil, and aged balsamic", 19.0, "vegetarian,gluten_free", true, nil},
				{"Appetizer", "Calamari Fritti", "Lightly fried squid with marinara and lemon aioli", 17.0, "", false, nil},
				{"Appetizer", "Carpaccio di Manzo", "Thinly sliced raw beef with arugula, capers, and shaved Parmigiano", 21.0, "gluten_free,raw", false, nil},
				{"Pasta", "Cacio e Pepe", "Handmade tonnarelli with Pecorino Romano and black pepper", 24.0, "vegetarian", true, intPtr(620)},
				{"Pasta", "Pappardelle Bolognese", "Wide ribbon pasta with slow-cooked beef and pork ragù", 26.0, "", false, intPtr(780)},
				{"Pasta", "Lobster Linguine", "Fresh linguine with butter-poached lobster, cherry tomatoes, and basil", 38.0, "", true, intPtr(710)},
				{"Pasta", "Gnocchi alla Sorrentina", "Potato gnocchi baked with tomato sauce, mozzarella, and basil", 22.0, "vegetarian", false, intPtr(580)},
				{"Main", "Branzino al Forno", "Whole roasted Mediterranean sea bass with lemon, herbs, and olive oil", 42.0, "gluten_free", false, intPtr(450)},
				{"Main", "Osso Buco", "Braised veal shank with saffron risotto and gremolata", 48.0, "gluten_free", false, intPtr(820)},
				{"Main", "Chicken Parmigiana", "Breaded chicken cutlet with marinara and melted mozzarella", 28.0, "", false, intPtr(750)},
				{"Pizza", "Margherita DOP", "San Marzano tomatoes, fior di latte, fresh basil, EVO", 20.0, "vegetarian", true, intPtr(680)},
				{"Pizza", "Truffle Funghi", "Wild mushrooms, truffle cream, fontina, and thyme", 26.0, "vegetarian", false, intPtr(720)},
				{"Dessert", "Tiramisu", "Classic layered mascarpone, espresso-soaked ladyfingers, cocoa", 14.0, "vegetarian", true, intPtr(480)},
				{"Dessert", "Panna Cotta", "Vanilla bean panna cotta with seasonal berry compote", 13.0, "vegetarian,gluten_free", false, intPtr(350)},
				{"Drink", "Aperol Spritz", "Aperol, prosecco, soda water", 16.0, "vegan,gluten_free", true, intPtr(180)},
				{"Drink", "Negroni", "Gin, Campari, sweet vermouth", 18.0, "vegan,gluten_free", false, intPtr(200)},
			},
		},
		// 2. Sakura House
		{
			Info: models.Restaurant{
				Name:        "Sakura House",
				Description: "Contemporary Japanese restaurant specializing in omakase sushi, handrolls, and seasonal izakaya plates. Fish flown in daily from Tsukiji Market.",
				Cuisines:    "Japanese,Sushi",
				PriceRange:  models.PriceUpscale,
				Address:     "88 E 10th St",
				City:        "New York",
				State:       "NY",
				ZipCode:     "10003",
				Country:     "US",
				Latitude:    ptr(40.7299),
				Longitude:   ptr(-73.9908),
				Phone:       "+1-212-555-0188",
				Email:       "hello@sakurahouse.example.com",
				Website:     "https://sakurahouse.example.com",
				Features:    "wifi,wheelchair_accessible",
				TotalSeats:  45,
				Rating:      ptr(4.8),
				ReviewCount: 567,
			},
			WeekdayOpen: "12:00", WeekdayClose: "22:30", WeekendOpen: "12:00", WeekendClose: "23:00",
			ClosedDays: []string{"monday"},
			Menu: []menuEntry{
				{"Sushi", "Salmon Nigiri (2pc)", "Fresh Atlantic salmon over seasoned rice", 8.0, "gluten_free,raw", false, intPtr(120)},
				{"Sushi", "Toro Nigiri (2pc)", "Fatty bluefin tuna belly", 18.0, "gluten_free,raw", true, intPtr(140)},
				{"Sushi", "Uni Nigiri (2pc)", "Santa Barbara sea urchin", 22.0, "gluten_free,raw", true, intPtr(130)},
				{"Sushi", "Hamachi Nigiri (2pc)", "Japanese yellowtail", 10.0, "gluten_free,raw", false, intPtr(110)},
				{"Rolls", "Dragon Roll", "Shrimp tempura, avocado, eel, eel sauce", 19.0, "", true, intPtr(380)},
				{"Rolls", "Spicy Tuna Roll", "Tuna, spicy mayo, cucumber, sesame", 16.0, "spicy,raw", false, intPtr(290)},
				{"Rolls", "Veggie Roll", "Avocado, cucumber, carrot, asparagus", 12.0, "vegan", false, intPtr(220)},
				{"Izakaya", "Edamame", "Steamed soybeans with sea salt", 7.0, "vegan,gluten_free", false, intPtr(120)},
				{"Izakaya", "Gyoza (6pc)", "Pan-fried pork dumplings with ponzu", 12.0, "", false, intPtr(320)},
				{"Izakaya", "Karaage", "Japanese fried chicken with kewpie mayo", 14.0, "", true, intPtr(450)},
				{"Izakaya", "Agedashi Tofu", "Fried silken tofu in dashi broth", 11.0, "vegetarian", false, intPtr(280)},
				{"Main", "Wagyu A5 Steak (4oz)", "Japanese A5 wagyu with wasabi and pink salt", 85.0, "gluten_free", true, intPtr(520)},
				{"Main", "Chirashi Bowl", "Chef's selection sashimi over sushi rice", 34.0, "gluten_free,raw", true, intPtr(480)},
				{"Main", "Ramen Tonkotsu", "Rich pork bone broth, chashu, soft egg, nori", 22.0, "", false, intPtr(680)},
				{"Dessert", "Mochi Ice Cream (3pc)", "Green tea, strawberry, and black sesame", 10.0, "vegetarian,gluten_free", false, intPtr(240)},
				{"Drink", "Sake Flight", "Three 2oz pours of premium sake", 24.0, "vegan,gluten_free", false, intPtr(180)},
				{"Drink", "Japanese Whisky Highball", "Suntory Toki with sparkling water", 16.0, "vegan,gluten_free", false, intPtr(140)},
			},
		},
		// 3. El Jardín
		{
			Info: models.Restaurant{
				Name:        "El Jardín",
				Description: "Vibrant Mexican restaurant with a lush garden patio. Traditional recipes from Oaxaca and Mexico City with a modern twist. Famous for our tableside guacamole and mezcal cocktails.",
				Cuisines:    "Mexican,Latin American",
				PriceRange:  models.PriceModerate,
				Address:     "2847 Mission St",
				City:        "San Francisco",
				State:       "CA",
				ZipCode:     "94110",
				Country:     "US",
				Latitude:    ptr(37.7520),
				Longitude:   ptr(-122.4189),
				Phone:       "+1-415-555-0284",
				Email:       "hola@eljardin.example.com",
				Website:     "https://eljardin.example.com",
				Features:    "outdoor_seating,live_music,parking,delivery,takeout,pet_friendly",
				TotalSeats:  120,
				Rating:      ptr(4.5),
				ReviewCount: 891,
			},
			Menu: []menuEntry{
				{"Appetizer", "Tableside Guacamole", "Made fresh at your table with avocado, lime, cilantro, jalapeño", 15.0, "vegan,gluten_free", true, intPtr(320)},
				{"Appetizer", "Elote", "Grilled street corn with cotija, mayo, chile, lime", 9.0, "vegetarian,gluten_free", true, intPtr(280)},
				{"Appetizer", "Queso Fundido", "Melted Oaxacan cheese with chorizo and warm tortillas", 14.0, "", false, intPtr(450)},
				{"Tacos", "Carnitas Tacos (3)", "Slow-roasted pork, pickled onion, salsa verde, cilantro", 16.0, "gluten_free", true, intPtr(520)},
				{"Tacos", "Al Pastor Tacos (3)", "Spit-roasted marinated pork with pineapple and onion", 16.0, "gluten_free", true, intPtr(490)},
				{"Tacos", "Fish Tacos (3)", "Beer-battered cod, cabbage slaw, chipotle crema", 17.0, "", false, intPtr(480)},
				{"Tacos", "Mushroom Tacos (3)", "Sautéed wild mushrooms, black beans, avocado, salsa roja", 15.0, "vegan,gluten_free", false, intPtr(380)},
				{"Main", "Mole Negro", "Chicken in traditional Oaxacan mole negro with 30+ ingredients", 28.0, "gluten_free", true, intPtr(680)},
				{"Main", "Enchiladas Suizas", "Chicken enchiladas with tomatillo cream sauce and queso fresco", 22.0, "", false, intPtr(590)},
				{"Main", "Carne Asada", "Grilled skirt steak with rice, beans, pico, guacamole", 26.0, "gluten_free", false, intPtr(720)},
				{"Main", "Chile Relleno", "Roasted poblano stuffed with cheese, walnut cream sauce", 20.0, "vegetarian,gluten_free", false, intPtr(520)},
				{"Dessert", "Churros con Chocolate", "Cinnamon sugar churros with Mexican hot chocolate dipping sauce", 12.0, "vegetarian", true, intPtr(420)},
				{"Dessert", "Tres Leches Cake", "Three-milk sponge cake with whipped cream and berries", 11.0, "vegetarian", false, intPtr(380)},
				{"Drink", "Mezcal Margarita", "Mezcal, fresh lime, agave, Tajín rim", 16.0, "vegan,gluten_free", true, intPtr(200)},
				{"Drink", "Horchata", "House-made rice milk with cinnamon and vanilla", 6.0, "vegan,gluten_free", false, intPtr(180)},
				{"Drink", "Agua de Jamaica", "Hibiscus iced tea with lime", 5.0, "vegan,gluten_free", false, intPtr(60)},
			},
		},
		// 4. The Green Plate
		{
			Info: models.Restaurant{
				Name:        "The Green Plate",
				Description: "100% plant-based restaurant serving creative vegan dishes made from locally sourced, organic ingredients. Zero-waste kitchen. Perfect for health-conscious diners.",
				Cuisines:    "Vegan,American,Health Food",
				PriceRange:  models.PriceModerate,
				Address:     "456 Abbot Kinney Blvd",
				City:        "Los Angeles",
				State:       "CA",
				ZipCode:     "90291",
				Country:     "US",
				Latitude:    ptr(33.9925),
				Longitude:   ptr(-118.4694),
				Phone:       "+1-310-555-0456",
				Email:       "eat@thegreenplate.example.com",
				Website:     "https://thegreenplate.example.com",
				Features:    "outdoor_seating,wifi,delivery,takeout,wheelchair_accessible,pet_friendly",
				TotalSeats:  60,
				Rating:      ptr(4.6),
				ReviewCount: 445,
			},
			WeekdayOpen: "08:00", WeekdayClose: "21:00", WeekendOpen: "08:00", WeekendClose: "22:00",
			Menu: []menuEntry{
				{"Breakfast", "Açaí Power Bowl", "Açaí, banana, granola, coconut, chia seeds, local berries", 16.0, "vegan,gluten_free,organic", true, intPtr(380)},
				{"Breakfast", "Avocado Toast", "Sourdough, smashed avocado, everything seasoning, microgreens, hemp seeds", 14.0, "vegan", true, intPtr(320)},
				{"Breakfast", "Tofu Scramble", "Seasoned tofu with black beans, peppers, mushrooms, salsa verde", 15.0, "vegan,gluten_free", false, intPtr(350)},
				{"Appetizer", "Cauliflower Wings", "Crispy battered cauliflower with buffalo or BBQ sauce", 13.0, "vegan", true, intPtr(280)},
				{"Appetizer", "Raw Zucchini Carpaccio", "Paper-thin zucchini with lemon, olive oil, pine nuts, mint", 12.0, "vegan,gluten_free,raw,organic", false, intPtr(180)},
				{"Main", "Impossible Burger", "Plant-based patty, vegan cheddar, lettuce, tomato, special sauce", 19.0, "vegan", true, intPtr(580)},
				{"Main", "Buddha Bowl", "Quinoa, roasted sweet potato, chickpeas, tahini, avocado, kale", 18.0, "vegan,gluten_free,organic", true, intPtr(520)},
				{"Main", "Mushroom Risotto", "Arborio rice with wild mushrooms, truffle oil, nutritional yeast", 22.0, "vegan,gluten_free", false, intPtr(480)},
				{"Main", "Jackfruit Tacos", "BBQ jackfruit, purple cabbage slaw, cashew crema, corn tortillas", 17.0, "vegan,gluten_free", false, intPtr(420)},
				{"Main", "Pad Thai", "Rice noodles, tofu, vegetables, tamarind sauce, crushed peanuts", 19.0, "vegan,gluten_free", false, intPtr(460)},
				{"Dessert", "Chocolate Avocado Mousse", "Rich dark chocolate mousse made with avocado, coconut cream", 11.0, "vegan,gluten_free,organic", true, intPtr(320)},
				{"Dessert", "Raw Cheesecake", "Cashew-based cheesecake with a date-walnut crust and berry coulis", 13.0, "vegan,gluten_free,raw", false, intPtr(380)},
				{"Drink", "Green Goddess Smoothie", "Kale, spinach, banana, mango, ginger, coconut water", 10.0, "vegan,gluten_free,organic", true, intPtr(220)},
				{"Drink", "Oat Milk Latte", "Espresso with house-made oat milk", 7.0, "vegan", false, intPtr(150)},
				{"Drink", "Cold-Pressed Juice Flight", "Four 4oz juices: beet-carrot, green, citrus, ginger-turmeric", 14.0, "vegan,gluten_free,raw,organic", false, intPtr(160)},
			},
		},
		// 5. Maison Laurent
		{
			Info: models.Restaurant{
				Name:        "Maison Laurent",
				Description: "Refined French fine dining with seasonal tasting menus and sommelier-curated wine pairings. A Michelin-starred experience in the heart of Chicago.",
				Cuisines:    "French,European",
				PriceRange:  models.PriceFineDining,
				Address:     "65 W Walton St",
				City:        "Chicago",
				State:       "IL",
				ZipCode:     "60610",
				Country:     "US",
				Latitude:    ptr(41.9005),
				Longitude:   ptr(-87.6290),
				Phone:       "+1-312-555-0065",
				Email:       "concierge@maisonlaurent.example.com",
				Website:     "https://maisonlaurent.example.com",
				Features:    "wifi,wheelchair_accessible,parking",
				TotalSeats:  40,
				Rating:      ptr(4.9),
				ReviewCount: 218,
			},
			WeekdayOpen: "17:30", WeekdayClose: "22:00", WeekendOpen: "17:00", WeekendClose: "22:30",
			ClosedDays: []string{"monday", "tuesday"},
			Menu: []menuEntry{
				{"Amuse-Bouche", "Foie Gras Bonbon", "Seared foie gras in a dark chocolate shell with fleur de sel", 0.0, "", false, intPtr(180)},
				{"Appetizer", "Tartare de Boeuf", "Hand-cut beef tartare, quail egg yolk, cornichon, dijon", 28.0, "gluten_free,raw", true, intPtr(280)},
				{"Appetizer", "Soupe à l'Oignon", "Classic French onion soup with Gruyère crouton", 18.0, "", false, intPtr(320)},
				{"Appetizer", "Escargot de Bourgogne", "Burgundy snails in garlic-herb butter", 24.0, "", false, intPtr(260)},
				{"Main", "Canard à l'Orange", "Roasted duck breast with orange gastrique, fondant potato, haricots verts", 52.0, "gluten_free", true, intPtr(680)},
				{"Main", "Filet de Boeuf", "Prime beef tenderloin, bordelaise sauce, pommes purée, seasonal vegetables", 62.0, "gluten_free", true, intPtr(720)},
				{"Main", "Dover Sole Meunière", "Whole Dover sole, brown butter, lemon, capers, parsley", 58.0, "gluten_free", false, intPtr(420)},
				{"Main", "Risotto aux Truffes", "Carnaroli rice, black truffle, aged Parmesan, truffle butter", 45.0, "vegetarian,gluten_free", false, intPtr(560)},
				{"Cheese", "French Cheese Board", "Selection of 5 artisanal French cheeses with honeycomb and walnut bread", 28.0, "vegetarian", false, intPtr(480)},
				{"Dessert", "Crème Brûlée", "Madagascar vanilla bean custard with caramelized sugar", 16.0, "vegetarian,gluten_free", true, intPtr(380)},
				{"Dessert", "Tarte Tatin", "Caramelized apple tart with crème fraîche", 18.0, "vegetarian", false, intPtr(420)},
				{"Dessert", "Soufflé au Chocolat", "Dark chocolate soufflé with vanilla crème anglaise (20 min)", 22.0, "vegetarian", true, intPtr(520)},
				{"Drink", "Sommelier's Wine Pairing", "5-course pairing selected by our sommelier", 95.0, "vegan,gluten_free", true, intPtr(400)},
				{"Drink", "Champagne Kir Royale", "Crème de cassis with vintage Champagne", 24.0, "vegan,gluten_free", false, intPtr(160)},
			},
		},
		// 6. Spice Route
		{
			Info: models.Restaurant{
				Name:        "Spice Route",
				Description: "Award-winning Indian restaurant featuring regional dishes from Kerala, Punjab, and Rajasthan. Tandoor oven, house-ground spice blends, and an extensive vegetarian menu.",
				Cuisines:    "Indian,South Asian",
				PriceRange:  models.PriceModerate,
				Address:     "1120 Connecticut Ave NW",
				City:        "Washington",
				State:       "DC",
				ZipCode:     "20036",
				Country:     "US",
				Latitude:    ptr(38.9060),
				Longitude:   ptr(-77.0408),
				Phone:       "+1-202-555-1120",
				Email:       "info@spiceroute.example.com",
				Website:     "https://spiceroute.example.com",
				Features:    "delivery,takeout,wifi,wheelchair_accessible,parking",
				TotalSeats:  90,
				Rating:      ptr(4.4),
				ReviewCount: 623,
			},
			Menu: []menuEntry{
				{"Appetizer", "Samosa (2pc)", "Crispy pastry filled with spiced potatoes and peas, tamarind chutney", 8.0, "vegetarian,vegan", true, intPtr(320)},
				{"Appetizer", "Chicken Tikka", "Tandoor-roasted chicken marinated in yogurt and spices", 14.0, "gluten_free", true, intPtr(280)},
				{"Appetizer", "Paneer Tikka", "Grilled cottage cheese with peppers and onions", 13.0, "vegetarian,gluten_free", false, intPtr(300)},
				{"Main", "Butter Chicken", "Tandoori chicken in creamy tomato-butter sauce", 22.0, "gluten_free", true, intPtr(620)},
				{"Main", "Lamb Rogan Josh", "Kashmiri slow-cooked lamb in aromatic red chili sauce", 24.0, "gluten_free,spicy", true, intPtr(580)},
				{"Main", "Palak Paneer", "Fresh spinach and cottage cheese in a mild spiced gravy", 18.0, "vegetarian,gluten_free", true, intPtr(380)},
				{"Main", "Chana Masala", "Chickpeas in a tangy tomato-onion gravy with cumin", 16.0, "vegan,gluten_free", false, intPtr(340)},
				{"Main", "Dal Makhani", "Creamy black lentils slow-cooked overnight with butter and cream", 17.0, "vegetarian,gluten_free", true, intPtr(420)},
				{"Main", "Shrimp Vindaloo", "Goan-style shrimp in a fiery vinegar-chili sauce", 24.0, "gluten_free,spicy", false, intPtr(380)},
				{"Bread", "Garlic Naan", "Leavened bread with garlic and butter from the tandoor", 5.0, "vegetarian", true, intPtr(260)},
				{"Bread", "Roti", "Whole wheat flatbread", 4.0, "vegan", false, intPtr(180)},
				{"Rice", "Biryani (Chicken)", "Basmati rice layered with spiced chicken, saffron, fried onions", 20.0, "gluten_free", true, intPtr(580)},
				{"Rice", "Vegetable Biryani", "Basmati rice with seasonal vegetables and aromatic spices", 17.0, "vegan,gluten_free", false, intPtr(440)},
				{"Dessert", "Gulab Jamun (3pc)", "Fried milk dumplings in rose-cardamom syrup", 9.0, "vegetarian", true, intPtr(360)},
				{"Dessert", "Mango Lassi", "Thick yogurt smoothie with Alphonso mango", 7.0, "vegetarian,gluten_free", false, intPtr(220)},
				{"Drink", "Masala Chai", "House-brewed spiced tea with milk", 5.0, "vegetarian,gluten_free", false, intPtr(80)},
			},
		},
		// 7. Burger & Barrel
		{
			Info: models.Restaurant{
				Name:        "Burger & Barrel",
				Description: "Classic American burger joint with craft beers on tap and hand-cut fries. Grass-fed beef, house-made buns, and creative toppings. The best burger in Austin.",
				Cuisines:    "American,Burgers",
				PriceRange:  models.PriceBudget,
				Address:     "1501 S Congress Ave",
				City:        "Austin",
				State:       "TX",
				ZipCode:     "78704",
				Country:     "US",
				Latitude:    ptr(30.2477),
				Longitude:   ptr(-97.7497),
				Phone:       "+1-512-555-1501",
				Email:       "eat@burgerbarrel.example.com",
				Website:     "https://burgerbarrel.example.com",
				Features:    "outdoor_seating,wifi,delivery,takeout,parking,pet_friendly",
				TotalSeats:  100,
				Rating:      ptr(4.3),
				ReviewCount: 1205,
			},
			WeekdayOpen: "11:00", WeekdayClose: "23:00", WeekendOpen: "10:00", WeekendClose: "00:00",
			Menu: []menuEntry{
				{"Burger", "Classic Smash Burger", "Double grass-fed beef patties, American cheese, lettuce, tomato, pickles, special sauce", 14.0, "", true, intPtr(780)},
				{"Burger", "BBQ Bacon Burger", "Beef patty, smoked bacon, cheddar, crispy onion rings, BBQ sauce", 16.0, "", true, intPtr(920)},
				{"Burger", "Mushroom Swiss Burger", "Beef patty, sautéed mushrooms, Swiss cheese, garlic aioli", 15.0, "", false, intPtr(820)},
				{"Burger", "Veggie Burger", "House-made black bean and quinoa patty with avocado and chipotle mayo", 13.0, "vegetarian", false, intPtr(560)},
				{"Burger", "The Texan", "Beef patty, fried egg, jalapeños, pepper jack, green chile salsa", 17.0, "spicy", true, intPtr(880)},
				{"Sides", "Hand-Cut Fries", "Crispy fries with sea salt", 6.0, "vegan,gluten_free", true, intPtr(380)},
				{"Sides", "Sweet Potato Fries", "With chipotle ketchup", 7.0, "vegan,gluten_free", false, intPtr(340)},
				{"Sides", "Onion Rings", "Beer-battered thick-cut onion rings", 8.0, "vegetarian", false, intPtr(420)},
				{"Sides", "Mac & Cheese", "Three-cheese baked mac with breadcrumb crust", 9.0, "vegetarian", true, intPtr(520)},
				{"Chicken", "Nashville Hot Chicken Sandwich", "Spicy fried chicken, pickles, coleslaw, brioche bun", 15.0, "spicy", true, intPtr(750)},
				{"Chicken", "Grilled Chicken Wrap", "Grilled chicken, lettuce, tomato, ranch, flour tortilla", 13.0, "", false, intPtr(520)},
				{"Dessert", "Milkshake", "Vanilla, chocolate, or strawberry — hand-spun", 8.0, "vegetarian,gluten_free", true, intPtr(580)},
				{"Dessert", "Apple Pie à la Mode", "Warm apple pie with vanilla ice cream", 9.0, "vegetarian", false, intPtr(480)},
				{"Drink", "Craft Beer Flight", "Four 5oz pours from our rotating taps", 14.0, "vegan", true, intPtr(300)},
				{"Drink", "House Lemonade", "Fresh-squeezed lemonade with mint", 5.0, "vegan,gluten_free", false, intPtr(120)},
			},
		},
		// 8. Jade Palace
		{
			Info: models.Restaurant{
				Name:        "Jade Palace",
				Description: "Elegant Cantonese and Sichuan restaurant with dim sum brunch, Peking duck, and a private dining room. Family-owned for three generations.",
				Cuisines:    "Chinese,Cantonese,Sichuan",
				PriceRange:  models.PriceModerate,
				Address:     "718 Jackson St",
				City:        "San Francisco",
				State:       "CA",
				ZipCode:     "94133",
				Country:     "US",
				Latitude:    ptr(37.7961),
				Longitude:   ptr(-122.4075),
				Phone:       "+1-415-555-0718",
				Email:       "info@jadepalace.example.com",
				Website:     "https://jadepalace.example.com",
				Features:    "wifi,wheelchair_accessible,parking,takeout,delivery",
				TotalSeats:  150,
				Rating:      ptr(4.5),
				ReviewCount: 734,
			},
			WeekdayOpen: "10:30", WeekdayClose: "22:00", WeekendOpen: "09:00", WeekendClose: "22:30",
			Menu: []menuEntry{
				{"Dim Sum", "Har Gow (4pc)", "Crystal shrimp dumplings", 8.0, "gluten_free", true, intPtr(160)},
				{"Dim Sum", "Siu Mai (4pc)", "Pork and shrimp dumplings", 8.0, "", true, intPtr(200)},
				{"Dim Sum", "Char Siu Bao (3pc)", "BBQ pork steamed buns", 7.0, "", true, intPtr(320)},
				{"Dim Sum", "Cheung Fun", "Rice noodle rolls with shrimp", 9.0, "gluten_free", false, intPtr(220)},
				{"Dim Sum", "Turnip Cake (3pc)", "Pan-fried radish cake with XO sauce", 7.0, "vegetarian", false, intPtr(240)},
				{"Main", "Peking Duck (whole)", "Carved tableside with pancakes, scallions, hoisin", 68.0, "", true, intPtr(880)},
				{"Main", "Kung Pao Chicken", "Diced chicken with peanuts, chilies, Sichuan peppercorn", 19.0, "spicy", true, intPtr(520)},
				{"Main", "Mapo Tofu", "Silken tofu in spicy fermented bean sauce with pork", 16.0, "spicy", true, intPtr(380)},
				{"Main", "Sweet & Sour Fish", "Crispy fish fillets in house-made sweet and sour sauce", 22.0, "", false, intPtr(480)},
				{"Main", "Beef Chow Fun", "Wide rice noodles with beef and bean sprouts", 18.0, "", false, intPtr(580)},
				{"Main", "Vegetable Fried Rice", "Wok-tossed rice with seasonal vegetables and egg", 14.0, "vegetarian", false, intPtr(440)},
				{"Soup", "Hot and Sour Soup", "Traditional soup with tofu, bamboo, mushrooms, egg", 10.0, "spicy", false, intPtr(180)},
				{"Soup", "Wonton Soup", "Pork wontons in clear chicken broth", 11.0, "", false, intPtr(220)},
				{"Dessert", "Mango Pudding", "Creamy mango pudding with coconut cream", 8.0, "vegetarian,gluten_free", true, intPtr(260)},
				{"Drink", "Jasmine Tea (pot)", "Premium jasmine green tea", 6.0, "vegan,gluten_free", false, intPtr(0)},
				{"Drink", "Lychee Martini", "Vodka, lychee liqueur, fresh lychee", 14.0, "vegan,gluten_free", false, intPtr(180)},
			},
		},
	}
}

func main() {
	cfg := config.Load()
	database.Init(cfg)

	// Check if already seeded
	var count int64
	database.DB.Model(&models.Restaurant{}).Count(&count)
	if count > 0 {
		fmt.Println("Database already has data — skipping seed. Delete agenteats.db to re-seed.")
		return
	}

	fmt.Println("🌱 Seeding AgentEats database...")

	seedData := getSeedData()
	for _, entry := range seedData {
		r := entry.Info
		r.ID = models.NewID()
		r.IsActive = true

		if err := database.DB.Create(&r).Error; err != nil {
			log.Fatalf("Failed to create restaurant %s: %v", r.Name, err)
		}

		// Hours
		hours := buildHours(r.ID, entry.WeekdayOpen, entry.WeekdayClose, entry.WeekendOpen, entry.WeekendClose, entry.ClosedDays)
		for _, h := range hours {
			database.DB.Create(&h)
		}

		// Menu
		for _, m := range entry.Menu {
			item := models.MenuItem{
				ID:            models.NewID(),
				RestaurantID:  r.ID,
				Category:      m.Category,
				Name:          m.Name,
				Description:   m.Description,
				Price:         m.Price,
				Currency:      "USD",
				DietaryLabels: m.Labels,
				IsAvailable:   true,
				IsPopular:     m.Popular,
				Calories:      m.Calories,
			}
			database.DB.Create(&item)
		}

		fmt.Printf("  ✓ %s (%s) — %d menu items\n", r.Name, r.City, len(entry.Menu))
	}

	// Sample reservations
	var restaurants []models.Restaurant
	database.DB.Limit(3).Find(&restaurants)

	sampleReservations := []models.Reservation{
		{ID: models.NewID(), RestaurantID: restaurants[0].ID, CustomerName: "Alice Johnson", CustomerEmail: "alice@example.com", CustomerPhone: "+1-555-0101", PartySize: 2, Date: "2026-02-20", Time: "19:00", Status: models.StatusConfirmed},
		{ID: models.NewID(), RestaurantID: restaurants[1].ID, CustomerName: "Bob Chen", CustomerEmail: "bob@example.com", CustomerPhone: "+1-555-0102", PartySize: 4, Date: "2026-02-20", Time: "20:00", Status: models.StatusConfirmed, SpecialRequests: "Window table please"},
		{ID: models.NewID(), RestaurantID: restaurants[2].ID, CustomerName: "Carol Davis", CustomerEmail: "carol@example.com", PartySize: 6, Date: "2026-02-21", Time: "19:30", Status: models.StatusConfirmed, SpecialRequests: "Birthday celebration — can you do a cake?"},
	}

	for _, res := range sampleReservations {
		database.DB.Create(&res)
	}

	fmt.Printf("\n✅ Seeded %d restaurants with menus and sample reservations.\n", len(seedData))
	fmt.Println("   Run `go run ./cmd/api` to start the REST API server.")
	fmt.Println("   Run `go run ./cmd/mcp` to start the MCP server for AI agents.")
}
