/**
 * Editorial SVG figures for the public README.
 * Source of truth for shapes: the sibling .mmd files.
 * Do not use reserved Mermaid ids (graph, click, end, subgraph).
 */
import { writeFileSync } from "node:fs";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";

const dir = dirname(fileURLToPath(import.meta.url));

const INK = "#1A1814";
const MUTED = "#5C564E";
const PAPER = "#F4F0E8";
const CARD = "#FFFCF7";
const FONT = "ui-sans-serif, system-ui, -apple-system, Segoe UI, Helvetica, Arial, sans-serif";

function svg(w, h, title, desc, inner) {
  return `<?xml version="1.0" encoding="UTF-8"?>
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 ${w} ${h}" width="${w}" height="${h}" role="img" aria-labelledby="title desc">
  <title id="title">${esc(title)}</title>
  <desc id="desc">${esc(desc)}</desc>
  <defs>
    <marker id="arrow" viewBox="0 0 10 10" refX="9" refY="5" markerWidth="7" markerHeight="7" orient="auto">
      <path d="M 0 0 L 10 5 L 0 10 z" fill="${INK}"/>
    </marker>
  </defs>
  <rect width="${w}" height="${h}" fill="${PAPER}"/>
  <rect x="0.75" y="0.75" width="${w - 1.5}" height="${h - 1.5}" fill="none" stroke="#D9D1C3" stroke-width="1.5"/>
  ${inner}
</svg>
`;
}

function esc(s) {
  return s.replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;");
}

function kicker(x, y, text) {
  return `<text x="${x}" y="${y}" fill="${MUTED}" font-size="11" font-weight="600" letter-spacing="2.6" font-family="${FONT}">${esc(text)}</text>`;
}

function box(x, y, w, h, label, sub) {
  const mid = x + w / 2;
  const labelY = sub ? y + h / 2 - 5 : y + h / 2 + 5;
  const subT = sub
    ? `<text x="${mid}" y="${y + h / 2 + 16}" text-anchor="middle" fill="${MUTED}" font-size="12" font-family="${FONT}">${esc(sub)}</text>`
    : "";
  return `<g>
  <rect x="${x}" y="${y}" width="${w}" height="${h}" rx="2" fill="${CARD}" stroke="${INK}" stroke-width="1.35"/>
  <text x="${mid}" y="${labelY}" text-anchor="middle" fill="${INK}" font-size="14" font-weight="600" font-family="${FONT}">${esc(label)}</text>
  ${subT}
</g>`;
}

function arrow(x1, y1, x2, y2) {
  return `<line x1="${x1}" y1="${y1}" x2="${x2}" y2="${y2}" stroke="${INK}" stroke-width="1.35" marker-end="url(#arrow)"/>`;
}

function poly(points) {
  return `<polyline points="${points}" fill="none" stroke="${INK}" stroke-width="1.35" marker-end="url(#arrow)"/>`;
}

function diamond(cx, cy, rw, rh, line1, line2) {
  const d = `M${cx},${cy - rh} L${cx + rw},${cy} L${cx},${cy + rh} L${cx - rw},${cy} Z`;
  const t2 = line2
    ? `<text x="${cx}" y="${cy + 14}" text-anchor="middle" fill="${CARD}" font-size="12" font-weight="600" font-family="${FONT}">${esc(line2)}</text>`
    : "";
  const y1 = line2 ? cy - 4 : cy + 5;
  return `<g>
  <path d="${d}" fill="${INK}" stroke="${INK}" stroke-width="1.35"/>
  <text x="${cx}" y="${y1}" text-anchor="middle" fill="${CARD}" font-size="12" font-weight="600" font-family="${FONT}">${esc(line1)}</text>
  ${t2}
</g>`;
}

const control = svg(
  880,
  420,
  "Orchestra control plane",
  "The human talks to one conductor. The conductor reads the capability graph, curated skills, executors, and private memory.",
  `
  ${kicker(36, 36, "CONTROL PLANE")}
  ${box(320, 56, 240, 52, "Human")}
  ${arrow(440, 108, 440, 142)}
  ${box(280, 148, 320, 64, "One conductor", "this session")}
  ${box(36, 268, 186, 92, "Capability graph", "13 routes")}
  ${box(244, 268, 186, 92, "Skills · MCP", "load by job")}
  ${box(452, 268, 186, 92, "Executors", "four hosts")}
  ${box(660, 268, 184, 92, "Private Brain", "your memory")}
  ${arrow(400, 212, 129, 268)}
  ${arrow(420, 212, 337, 268)}
  ${arrow(460, 212, 545, 268)}
  ${arrow(480, 212, 752, 268)}
`
);

const loop = svg(
  880,
  400,
  "Orchestra eight-stage loop",
  "Discover, classify, research, synthesize, then a human gate. Approved work is implemented, verified on the real app, and remembered.",
  `
  ${kicker(36, 36, "THE LOOP")}
  ${box(36, 56, 188, 56, "1  Discover")}
  ${box(244, 56, 188, 56, "2  Classify")}
  ${box(452, 56, 188, 56, "3  Research")}
  ${box(660, 56, 184, 56, "4  Synthesize")}
  ${arrow(224, 84, 244, 84)}
  ${arrow(432, 84, 452, 84)}
  ${arrow(640, 84, 660, 84)}
  ${poly("752,112 752,148 110,148 110,196")}
  ${box(36, 196, 148, 56, "5  Design")}
  ${diamond(318, 224, 78, 44, "Human", "gate")}
  ${box(452, 196, 188, 56, "6  Implement")}
  ${box(660, 196, 184, 56, "7  Verify")}
  ${arrow(184, 224, 240, 224)}
  ${arrow(396, 224, 452, 224)}
  ${arrow(640, 224, 660, 224)}
  ${poly("752,252 752,292 442,292 442,312")}
  ${box(244, 312, 396, 52, "8  Remember what actually ran")}
  <text x="546" y="176" text-anchor="middle" fill="${MUTED}" font-size="11" font-family="${FONT}">reject returns to research</text>
`
);

const hosts = svg(
  880,
  320,
  "Four hosts, one contract",
  "orchestra-workflow syncs the same allowlisted skills to Cursor, Antigravity, Claude Code, and jcode.",
  `
  ${kicker(36, 36, "FOUR HOSTS")}
  ${box(36, 120, 220, 80, "orchestra-workflow", "contract · graph · skills")}
  ${arrow(256, 160, 292, 160)}
  ${box(292, 120, 148, 80, "sync-ides", "allowlist")}
  ${box(516, 48, 328, 48, "Cursor")}
  ${box(516, 112, 328, 48, "Antigravity")}
  ${box(516, 176, 328, 48, "Claude Code")}
  ${box(516, 240, 328, 48, "jcode")}
  ${arrow(440, 148, 516, 72)}
  ${arrow(440, 160, 516, 136)}
  ${arrow(440, 172, 516, 200)}
  ${arrow(440, 184, 516, 264)}
`
);

writeFileSync(join(dir, "control-plane.svg"), control);
writeFileSync(join(dir, "the-loop.svg"), loop);
writeFileSync(join(dir, "four-hosts.svg"), hosts);
console.log("wrote 3 svg");
