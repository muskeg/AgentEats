// AgentEats — Dashboard logic

document.addEventListener('DOMContentLoaded', () => {
  if (!UI.requireAuth()) return;
  initHoursGrid();
  loadMyRestaurants();
});

// --- State ---
let myRestaurants = [];

// --- Navigation ---

function showPanel(panelId) {
  // Hide all panels
  document.querySelectorAll('.panel').forEach(p => p.classList.remove('active'));
  document.getElementById('panel-' + panelId)?.classList.add('active');

  // Update sidebar
  document.querySelectorAll('.sidebar-nav a[data-panel]').forEach(a => {
    a.classList.toggle('active', a.dataset.panel === panelId);
  });

  // Close mobile sidebar
  document.querySelector('.sidebar')?.classList.remove('open');

  // Panel-specific init
  if (panelId === 'manage-menu') {
    populateMenuRestaurantPicker();
  }
}

function switchTab(tabId, btn) {
  // Deactivate all
  btn.closest('.tabs').querySelectorAll('.tab-btn').forEach(b => b.classList.remove('active'));
  const parent = btn.closest('.panel') || btn.closest('#menu-content');
  parent.querySelectorAll('.tab-content').forEach(t => t.classList.remove('active'));

  // Activate
  btn.classList.add('active');
  document.getElementById('tab-' + tabId)?.classList.add('active');
}

// --- Logout ---

function handleLogout() {
  AgentEatsAPI.logout();
  window.location.href = 'login.html';
}

// --- Rotate Key ---

async function handleRotateKey() {
  if (!confirm('Are you sure? Your current API key will stop working immediately. You will receive a new key.')) {
    return;
  }

  try {
    const result = await AgentEatsAPI.rotateApiKey();
    AgentEatsAPI.login(result.api_key);

    // Show the new key
    const msg = `New API key:\n\n${result.api_key}\n\nThis is the only time it will be shown. Copy it now!`;
    alert(msg);

    // Also try to copy
    UI.copyToClipboard(result.api_key);
    UI.showToast('API key rotated. New key copied to clipboard.', 'success');
  } catch (err) {
    UI.showToast(err.message, 'error');
  }
}

// --- Load Restaurants ---

async function loadMyRestaurants() {
  const container = document.getElementById('restaurant-list');

  try {
    const results = await AgentEatsAPI.listMyRestaurants();
    myRestaurants = Array.isArray(results) ? results : [];

    if (myRestaurants.length === 0) {
      container.innerHTML = `
        <div class="alert alert-info">
          No restaurants found. <a href="#" onclick="showPanel('add-restaurant'); return false;">Add your first restaurant</a> to get started.
        </div>
      `;
      return;
    }

    container.innerHTML = myRestaurants.map(r => `
      <div class="restaurant-item" onclick="editRestaurant('${r.id}')">
        <div>
          <h3>${UI.escapeHTML(r.name)}</h3>
          <span class="meta">${UI.escapeHTML(r.city || '')} · ${UI.escapeHTML(r.price_range || '')} · ${(r.cuisines || []).join(', ')}</span>
        </div>
        <span style="color: var(--color-text-muted); font-size: 0.75rem;">Edit →</span>
      </div>
    `).join('');

  } catch (err) {
    container.innerHTML = `<div class="alert alert-error">${UI.escapeHTML(err.message)}</div>`;
  }
}

// --- Edit Restaurant ---

async function editRestaurant(id) {
  showPanel('add-restaurant');

  // Update form title
  document.getElementById('restaurant-form-title').textContent = 'Edit Restaurant';
  document.getElementById('restaurant-form-desc').textContent = 'Update your restaurant details.';
  document.getElementById('restaurant-submit-btn').textContent = 'Save Changes';
  document.getElementById('edit-restaurant-id').value = id;

  try {
    const r = await AgentEatsAPI.getRestaurant(id);

    document.getElementById('r-name').value = r.name || '';
    document.getElementById('r-description').value = r.description || '';
    document.getElementById('r-cuisines').value = (r.cuisines || []).join(', ');
    document.getElementById('r-price-range').value = r.price_range || '$$';
    document.getElementById('r-address').value = r.address || '';
    document.getElementById('r-city').value = r.city || '';
    document.getElementById('r-state').value = r.state || '';
    document.getElementById('r-zip').value = r.zip_code || '';
    document.getElementById('r-country').value = r.country || 'US';
    document.getElementById('r-phone').value = r.phone || '';
    document.getElementById('r-email').value = r.email || '';
    document.getElementById('r-website').value = r.website || '';
    document.getElementById('r-seats').value = r.total_seats || 50;

    // Features
    const features = r.features || [];
    document.querySelectorAll('.features-grid input[type="checkbox"]').forEach(cb => {
      cb.checked = features.includes(cb.value);
    });

    // Hours
    if (r.hours && r.hours.length > 0) {
      r.hours.forEach(h => {
        const day = h.day.toLowerCase();
        const openInput = document.querySelector(`#hours-${day}-open`);
        const closeInput = document.querySelector(`#hours-${day}-close`);
        const closedCb = document.querySelector(`#hours-${day}-closed`);
        if (openInput) openInput.value = h.open_time || '09:00';
        if (closeInput) closeInput.value = h.close_time || '22:00';
        if (closedCb) closedCb.checked = h.is_closed || false;
      });
    }

  } catch (err) {
    UI.showToast('Failed to load restaurant: ' + err.message, 'error');
  }
}

// --- Hours Grid ---

function initHoursGrid() {
  const days = ['monday', 'tuesday', 'wednesday', 'thursday', 'friday', 'saturday', 'sunday'];
  const grid = document.getElementById('hours-grid');
  if (!grid) return;

  grid.innerHTML = days.map(day => `
    <div class="hours-row">
      <label>${day}</label>
      <input type="time" id="hours-${day}-open" class="form-input" value="09:00">
      <input type="time" id="hours-${day}-close" class="form-input" value="22:00">
      <label class="form-checkbox" style="font-size: 0.75rem;">
        <input type="checkbox" id="hours-${day}-closed"> Closed
      </label>
    </div>
  `).join('');
}

// --- Submit Restaurant ---

async function handleRestaurantSubmit(e) {
  e.preventDefault();

  const errorDiv = document.getElementById('restaurant-form-error');
  const submitBtn = document.getElementById('restaurant-submit-btn');
  errorDiv.style.display = 'none';

  // Gather features
  const features = [];
  document.querySelectorAll('.features-grid input[type="checkbox"]:checked').forEach(cb => {
    features.push(cb.value);
  });

  // Gather hours
  const days = ['monday', 'tuesday', 'wednesday', 'thursday', 'friday', 'saturday', 'sunday'];
  const hours = days.map(day => ({
    day,
    open_time: document.getElementById(`hours-${day}-open`).value || '09:00',
    close_time: document.getElementById(`hours-${day}-close`).value || '22:00',
    is_closed: document.getElementById(`hours-${day}-closed`).checked,
  }));

  const data = {
    name: document.getElementById('r-name').value.trim(),
    description: document.getElementById('r-description').value.trim(),
    cuisines: document.getElementById('r-cuisines').value.split(',').map(s => s.trim()).filter(Boolean),
    price_range: document.getElementById('r-price-range').value,
    address: document.getElementById('r-address').value.trim(),
    city: document.getElementById('r-city').value.trim(),
    state: document.getElementById('r-state').value.trim(),
    zip_code: document.getElementById('r-zip').value.trim(),
    country: document.getElementById('r-country').value.trim() || 'US',
    phone: document.getElementById('r-phone').value.trim(),
    email: document.getElementById('r-email').value.trim(),
    website: document.getElementById('r-website').value.trim(),
    features,
    total_seats: parseInt(document.getElementById('r-seats').value) || 50,
    hours,
  };

  UI.setLoading(submitBtn, true);

  try {
    const editId = document.getElementById('edit-restaurant-id').value;

    if (editId) {
      await AgentEatsAPI.updateRestaurant(editId, data);
      UI.showToast('Restaurant updated!', 'success');
    } else {
      await AgentEatsAPI.createRestaurant(data);
      UI.showToast('Restaurant created!', 'success');
    }

    // Reset and go back to list
    resetRestaurantForm();
    showPanel('my-restaurants');
    loadMyRestaurants();

  } catch (err) {
    errorDiv.textContent = err.message;
    errorDiv.style.display = 'block';
  } finally {
    UI.setLoading(submitBtn, false);
  }
}

function resetRestaurantForm() {
  document.getElementById('restaurant-form').reset();
  document.getElementById('edit-restaurant-id').value = '';
  document.getElementById('restaurant-form-title').textContent = 'Add Restaurant';
  document.getElementById('restaurant-form-desc').textContent = 'Fill in your restaurant details to get listed on AgentEats.';
  document.getElementById('restaurant-submit-btn').textContent = 'Create Restaurant';
  document.getElementById('restaurant-form-error').style.display = 'none';
  document.getElementById('r-country').value = 'US';
  document.getElementById('r-seats').value = '50';
  initHoursGrid();
}

// --- Menu Management ---

function populateMenuRestaurantPicker() {
  const select = document.getElementById('menu-restaurant-select');
  const existingVal = select.value;

  // Keep first option
  select.innerHTML = '<option value="">— Choose a restaurant —</option>';

  myRestaurants.forEach(r => {
    const opt = document.createElement('option');
    opt.value = r.id;
    opt.textContent = r.name;
    select.appendChild(opt);
  });

  if (existingVal) select.value = existingVal;
}

async function handleMenuRestaurantChange() {
  const restaurantId = document.getElementById('menu-restaurant-select').value;
  const content = document.getElementById('menu-content');

  if (!restaurantId) {
    content.style.display = 'none';
    return;
  }

  content.style.display = 'block';
  await loadCurrentMenu(restaurantId);
}

async function loadCurrentMenu(restaurantId) {
  const container = document.getElementById('current-menu-list');

  try {
    const menu = await AgentEatsAPI.getMenu(restaurantId);
    const categories = menu.categories || {};
    const categoryNames = Object.keys(categories);

    if (categoryNames.length === 0) {
      container.innerHTML = '<div class="alert alert-info">No menu items yet. Add some using the tabs above.</div>';
      return;
    }

    let html = '';
    categoryNames.forEach(cat => {
      html += `<h3 style="margin: var(--space-lg) 0 var(--space-sm); font-size: 1rem;">${UI.escapeHTML(cat)}</h3>`;
      categories[cat].forEach(item => {
        const dietary = (item.dietary_labels || []).join(', ');
        html += `
          <div class="menu-item-card">
            <div>
              <div class="name">${UI.escapeHTML(item.name)}${item.is_popular ? ' ⭐' : ''}${!item.is_available ? ' <span style="color: var(--color-error);">(unavailable)</span>' : ''}</div>
              <div class="details">${UI.escapeHTML(item.description || '')}${dietary ? ' · ' + UI.escapeHTML(dietary) : ''}</div>
            </div>
            <div class="price">${item.currency || 'USD'} ${item.price?.toFixed(2) || '0.00'}</div>
          </div>
        `;
      });
    });

    container.innerHTML = html;

  } catch (err) {
    container.innerHTML = `<div class="alert alert-error">${UI.escapeHTML(err.message)}</div>`;
  }
}

async function handleAddMenuItem(e) {
  e.preventDefault();

  const restaurantId = document.getElementById('menu-restaurant-select').value;
  if (!restaurantId) {
    UI.showToast('Please select a restaurant first.', 'warning');
    return;
  }

  const errorDiv = document.getElementById('menu-item-error');
  const submitBtn = document.getElementById('menu-item-submit');
  errorDiv.style.display = 'none';

  const item = {
    name: document.getElementById('mi-name').value.trim(),
    category: document.getElementById('mi-category').value.trim(),
    description: document.getElementById('mi-description').value.trim(),
    price: parseFloat(document.getElementById('mi-price').value),
    currency: document.getElementById('mi-currency').value.trim() || 'USD',
    dietary_labels: document.getElementById('mi-dietary').value.split(',').map(s => s.trim()).filter(Boolean),
    is_available: document.getElementById('mi-available').checked,
    is_popular: document.getElementById('mi-popular').checked,
  };

  UI.setLoading(submitBtn, true);

  try {
    await AgentEatsAPI.addMenuItem(restaurantId, item);
    UI.showToast('Menu item added!', 'success');
    document.getElementById('menu-item-form').reset();
    document.getElementById('mi-currency').value = 'USD';
    document.getElementById('mi-available').checked = true;
    await loadCurrentMenu(restaurantId);
    // Switch to current menu tab
    switchTab('current-menu', document.querySelector('.tab-btn'));
  } catch (err) {
    errorDiv.textContent = err.message;
    errorDiv.style.display = 'block';
  } finally {
    UI.setLoading(submitBtn, false);
  }
}

async function handleBulkImport() {
  const restaurantId = document.getElementById('menu-restaurant-select').value;
  if (!restaurantId) {
    UI.showToast('Please select a restaurant first.', 'warning');
    return;
  }

  const errorDiv = document.getElementById('import-error');
  errorDiv.style.display = 'none';

  const jsonText = document.getElementById('import-json').value.trim();
  if (!jsonText) {
    errorDiv.textContent = 'Please paste your menu items JSON.';
    errorDiv.style.display = 'block';
    return;
  }

  let items;
  try {
    items = JSON.parse(jsonText);
    if (!Array.isArray(items)) throw new Error('Expected a JSON array');
  } catch (err) {
    errorDiv.textContent = 'Invalid JSON: ' + err.message;
    errorDiv.style.display = 'block';
    return;
  }

  const strategy = document.getElementById('import-strategy').value;

  try {
    const result = await AgentEatsAPI.bulkImportMenu(restaurantId, { strategy, items });
    UI.showToast(`Imported ${result.imported} items (${result.strategy})!`, 'success');
    document.getElementById('import-json').value = '';
    await loadCurrentMenu(restaurantId);
    switchTab('current-menu', document.querySelector('.tab-btn'));
  } catch (err) {
    errorDiv.textContent = err.message;
    errorDiv.style.display = 'block';
  }
}
