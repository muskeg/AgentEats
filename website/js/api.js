// AgentEats API client
// All calls to the backend go through this module.

const API_BASE = window.AGENTEATS_API_BASE || 'https://agenteats.fly.dev';

// --- Low-level helpers ---

async function apiFetch(path, options = {}) {
  const url = `${API_BASE}${path}`;
  const headers = { 'Content-Type': 'application/json', ...options.headers };

  // Attach auth if available
  const apiKey = sessionStorage.getItem('ae_api_key');
  if (apiKey) {
    headers['Authorization'] = `Bearer ${apiKey}`;
  }

  const res = await fetch(url, { ...options, headers });

  if (res.status === 429) {
    throw new ApiError('Too many requests — please wait a moment and try again.', 429);
  }

  const data = await res.json().catch(() => null);

  if (!res.ok) {
    const msg = data?.error || `Request failed (${res.status})`;
    throw new ApiError(msg, res.status);
  }

  return data;
}

class ApiError extends Error {
  constructor(message, status) {
    super(message);
    this.name = 'ApiError';
    this.status = status;
  }
}

// --- Auth ---

async function registerOwner(name, email) {
  return apiFetch('/owners/register', {
    method: 'POST',
    body: JSON.stringify({ name, email }),
  });
}

async function rotateApiKey() {
  return apiFetch('/owners/rotate-key', {
    method: 'POST',
  });
}

// Validate the stored API key by hitting an authenticated endpoint.
// We use rotate-key with a GET-like check — but since there's no
// dedicated "whoami" endpoint, we'll try a lightweight approach:
// attempt to list own restaurants. For now we just check that sessionStorage has a key.
function isLoggedIn() {
  return !!sessionStorage.getItem('ae_api_key');
}

function login(apiKey) {
  sessionStorage.setItem('ae_api_key', apiKey.trim());
}

function logout() {
  sessionStorage.removeItem('ae_api_key');
}

// --- Restaurants ---

async function createRestaurant(data) {
  return apiFetch('/restaurants', {
    method: 'POST',
    body: JSON.stringify(data),
  });
}

async function updateRestaurant(id, data) {
  return apiFetch(`/restaurants/${id}`, {
    method: 'PUT',
    body: JSON.stringify(data),
  });
}

async function getRestaurant(id) {
  return apiFetch(`/restaurants/${id}`);
}

async function searchRestaurants(params = {}) {
  const qs = new URLSearchParams(params).toString();
  return apiFetch(`/restaurants${qs ? '?' + qs : ''}`);
}

async function listMyRestaurants() {
  return apiFetch('/owners/restaurants');
}

// --- Menu ---

async function addMenuItem(restaurantId, item) {
  return apiFetch(`/restaurants/${restaurantId}/menu/items`, {
    method: 'POST',
    body: JSON.stringify(item),
  });
}

async function bulkImportMenu(restaurantId, data) {
  return apiFetch(`/restaurants/${restaurantId}/menu/import`, {
    method: 'POST',
    body: JSON.stringify(data),
  });
}

async function getMenu(restaurantId) {
  return apiFetch(`/restaurants/${restaurantId}/menu`);
}

// Export for use in other scripts
window.AgentEatsAPI = {
  ApiError,
  registerOwner,
  rotateApiKey,
  isLoggedIn,
  login,
  logout,
  createRestaurant,
  updateRestaurant,
  getRestaurant,
  searchRestaurants,
  listMyRestaurants,
  addMenuItem,
  bulkImportMenu,
  getMenu,
};
