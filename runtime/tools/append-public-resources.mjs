// Append unverified public reference rows. Does not set last_verified.
import fs from 'node:fs';
import path from 'node:path';

const file = path.join(import.meta.dirname, '..', '..', 'registries', 'resources.json');
const rows = JSON.parse(fs.readFileSync(file, 'utf8'));
const have = new Set(rows.map((row) => row.id));

const add = [
  ['godly', 'Godly', 'https://godly.website', 'component_gallery', 'Sites with a strong first screen.'],
  ['awwwards', 'Awwwards', 'https://www.awwwards.com', 'web_reference', 'Award galleries. Study the principle, do not copy the site.'],
  ['siteinspire', 'SiteInspire', 'https://www.siteinspire.com', 'web_reference', 'Editorial layout references.'],
  ['lapa-ninja', 'Lapa Ninja', 'https://www.lapa.ninja', 'web_reference', 'Landing page structure references.'],
  ['mobbin', 'Mobbin', 'https://mobbin.com', 'web_reference', 'Product UI flows.'],
  ['refero-design-gallery', 'Refero', 'https://refero.design', 'web_reference', 'Product screens. The refero-design skill stays the method.'],
  ['codrops', 'Codrops', 'https://tympanus.net/codrops/', 'web_reference', 'Interaction articles and demos.'],
  ['fontshare', 'Fontshare', 'https://www.fontshare.com', 'asset_archive', 'Free display and text faces.'],
  ['fontsource', 'Fontsource', 'https://fontsource.org', 'npm_package', 'Self-hosted font packages.'],
  ['lucide', 'Lucide', 'https://lucide.dev', 'component_gallery', 'Stroke icons.'],
  ['phosphor-icons', 'Phosphor', 'https://phosphoricons.com', 'component_gallery', 'Icon set with weights.'],
  ['unsplash', 'Unsplash', 'https://unsplash.com', 'asset_archive', 'Photographs with a stated license on the photo.'],
  ['pexels', 'Pexels', 'https://www.pexels.com', 'asset_archive', 'Photographs and video with a stated license.'],
  ['rive', 'Rive', 'https://rive.app', 'web_reference', 'State-machine motion. Not a second scroll runtime.'],
  ['motion-dev', 'Motion', 'https://motion.dev', 'web_reference', 'Motion docs. GSAP still owns scroll on an Orchestra page.'],
  ['threejs', 'Three.js', 'https://threejs.org', 'web_reference', 'The renderer under ow-3d-engineering.'],
];

let added = 0;
for (const [id, name, url, source, rationale] of add) {
  if (have.has(id)) continue;
  rows.push({
    id,
    name,
    canonical_url: url,
    source_type: source,
    license: 'unspecified',
    category: ['reference'],
    representation: 'reference',
    routing_tags: ['reference', 'visual'],
    acquisition_method: 'link',
    runtime_method: 'lookup',
    status: 'ACTIVE',
    trigger_conditions: ['a visual task needs a public reference'],
    skip_conditions: ['copying the asset, the copy, or the logo'],
    rationale,
    usage_evidence: 'none recorded in this registry',
    lifecycle_stage: 'discovered',
    admission_state: 'discovered',
  });
  added += 1;
}

fs.writeFileSync(file, JSON.stringify(rows, null, 2) + '\n');
const verified = rows.filter((row) => String(row.last_verified || '').trim()).length;
console.log(`added ${added} total ${rows.length} verified ${verified} unverified ${rows.length - verified}`);
