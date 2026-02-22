// Load and render markdown documentation
async function loadDocs(mdPath) {
  const loadingEl = document.getElementById('docs-loading');
  const contentEl = document.getElementById('docs-content');

  // Generate GitHub-style heading slugs so TOC anchor links work
  const renderer = new marked.Renderer();
  const slugCounts = {};
  renderer.heading = function ({ text, depth }) {
    // Strip any inline HTML tags to get plain text for the slug
    const raw = text.replace(/<[^>]+>/g, '');
    let slug = raw
      .toLowerCase()
      .replace(/[^\w\s-]/g, '')   // remove non-word chars except spaces/hyphens
      .replace(/\s+/g, '-')       // spaces to hyphens
      .replace(/-+/g, '-')        // collapse multiple hyphens
      .replace(/^-|-$/g, '');     // trim leading/trailing hyphens

    // Handle duplicate slugs (append -1, -2, etc.)
    if (slugCounts[slug] !== undefined) {
      slugCounts[slug]++;
      slug = `${slug}-${slugCounts[slug]}`;
    } else {
      slugCounts[slug] = 0;
    }

    return `<h${depth} id="${slug}">${text}</h${depth}>`;
  };

  marked.setOptions({ renderer });

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

    // Handle in-page anchor clicks with smooth scroll
    contentEl.addEventListener('click', function (e) {
      const link = e.target.closest('a[href^="#"]');
      if (!link) return;
      e.preventDefault();
      const id = link.getAttribute('href').slice(1);
      const target = document.getElementById(id);
      if (target) {
        target.scrollIntoView({ behavior: 'smooth' });
        history.replaceState(null, '', '#' + id);
      }
    });
  } catch (err) {
    console.error(err);
    loadingEl.innerHTML = '<p>Failed to load documentation. Please try refreshing the page.</p>';
  }
}
