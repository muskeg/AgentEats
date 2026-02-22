// Load and render markdown documentation
async function loadDocs(mdPath) {
  const loadingEl = document.getElementById('docs-loading');
  const contentEl = document.getElementById('docs-content');

  try {
    const resp = await fetch(mdPath);
    if (!resp.ok) throw new Error(`Failed to load ${mdPath}`);

    let md = await resp.text();

    // Strip any GitHub release links (repo is closed-source)
    md = md.replace(/\[GitHub Releases\]\([^)]+\)/g, 'the releases page');

    // Render markdown to HTML
    contentEl.innerHTML = marked.parse(md);
    contentEl.style.display = '';
    loadingEl.style.display = 'none';

    // Smooth-scroll to hash if present
    if (window.location.hash) {
      const target = document.getElementById(window.location.hash.slice(1));
      if (target) target.scrollIntoView({ behavior: 'smooth' });
    }
  } catch (err) {
    console.error(err);
    loadingEl.innerHTML = '<p>Failed to load documentation. Please try refreshing the page.</p>';
  }
}
