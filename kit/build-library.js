const fs = require("fs");
const path = require("path");
const vault = process.env.ORCHESTRA_VAULT;
if (!vault) throw new Error("Set ORCHESTRA_VAULT");
const root = path.resolve(__dirname, "..");
const rows = require(path.join(root, "registries", "resources.json"));
const list = Array.isArray(rows) ? rows : rows.resources;
const lib = path.join(vault, "library");
const notes = path.join(lib, "resources");
fs.mkdirSync(notes, { recursive: true });
function q(v) {
  return '"' + String(v == null ? "" : v).replace(/\\/g, "\\\\").replace(/"/g, '\\"').replace(/\r?\n/g, " ") + '"';
}
function tagOf(status) {
  if (status === "ACTIVE") return "job";
  if (status === "REJECTED") return "refused";
  if (status === "PROJECT_ONLY") return "project";
  if (status === "EXPERIMENTAL") return "pilot";
  if (status === "PURCHASE_PENDING") return "shelf";
  return "shelf";
}
function fieldOf(row) {
  const blob = (Array.isArray(row.category) ? row.category.join(" ") : String(row.category || "")).toLowerCase();
  const rules = [
    ["3d", /3d|shader|webgl|r3f|three/],
    ["motion", /motion|animat|gsap/],
    ["mobile", /mobile|expo|android|ios/],
    ["docs", /doc|slide|paper|diagram|present/],
    ["security", /secur|pentest/],
    ["research", /research|scrap|crawl|search/],
    ["design", /design|ui|ux|visual|frontend|css|stitch/],
    ["data", /data|database|postgres|supabase/],
    ["agentic", /agent|skill|mcp|harness/]
  ];
  for (const [name, re] of rules) if (re.test(blob)) return name;
  return "general";
}
const hubs = {};
const refused = [];
for (const row of list) {
  const id = String(row.id || "untitled").replace(/[^a-z0-9-]+/gi, "-").toLowerCase();
  const field = fieldOf(row);
  const tag = tagOf(row.status);
  const url = row.canonical_url || row.documentation_url || "";
  const load = tag === "job" ? "Read the matching skill if one is installed. Otherwise open the URL." : tag === "refused" ? "Do not load." : "Open the URL when a task names it.";
  const body = [
    "---",
    "title: " + q(row.name || id),
    "field: " + q(field),
    "tag: " + q(tag),
    "status: " + q(row.status || ""),
    "url: " + q(url),
    "load: " + q(load),
    "---",
    "",
    "# " + (row.name || id),
    "",
    (row.rationale || "No rationale recorded.").replace(/\r?\n/g, " "),
    "",
    "| | |",
    "| --- | --- |",
    "| Use when | " + ((row.trigger_conditions || []).slice(0, 3).join("; ") || "A task names it.") + " |",
    "| Avoid | " + ((row.avoid_conditions || []).slice(0, 3).join("; ") || "") + " |",
    "| Load | " + load + " |",
    ""
  ].join("\n");
  fs.writeFileSync(path.join(notes, id + ".md"), body);
  (hubs[field] = hubs[field] || []).push({ id, name: row.name || id, tag, status: row.status || "" });
  if (tag === "refused") refused.push({ id, name: row.name || id, why: (row.rationale || "").replace(/\r?\n/g, " ") });
}
for (const [field, items] of Object.entries(hubs)) {
  const lines = ["---", "title: " + q(field), "tags: [library]", "---", "", "# " + field, "", "| Name | Tag | Status |", "| --- | --- | --- |"];
  for (const item of items.sort((a, b) => a.name.localeCompare(b.name))) {
    lines.push("| [[" + item.id + "]] | " + item.tag + " | " + item.status + " |");
  }
  lines.push("");
  fs.writeFileSync(path.join(lib, field + ".md"), lines.join("\n"));
}
const ref = ["---", "title: Refused", "---", "", "# Refused", "", "| Name | Why |", "| --- | --- |"];
for (const item of refused.sort((a, b) => a.name.localeCompare(b.name))) {
  ref.push("| [[" + item.id + "]] | " + item.why.slice(0, 180).replace(/\|/g, "/") + " |");
}
ref.push("");
fs.writeFileSync(path.join(lib, "refused.md"), ref.join("\n"));
const readme = [
  "---",
  "title: Library",
  "---",
  "",
  "# Library",
  "",
  "> [!summary]",
  "> Generated from the workflow registry. One note per resource.",
  "> The written verdicts for the zip shelf stay in the 3.3.5 resource ledger.",
  "",
  "| Tag | You say | Agent does |",
  "| --- | --- | --- |",
  "| job | use the name | Reads the skill |",
  "| shelf | use the name from library | Opens the note and the URL |",
  "| project | use it in this repo | Installs only in that repo |",
  "| pilot | name the pilot | Does not make it the default |",
  "| refused | | Declines and cites the note |",
  "",
  "Hubs are the other files in this folder. Notes are in `resources/`.",
  ""
].join("\n");
fs.writeFileSync(path.join(lib, "README.md"), readme);
const base = [
  "filters:",
  "  and:",
  "    - file.inFolder(\"library/resources\")",
  "properties:",
  "  tag:",
  "    displayName: Tag",
  "  status:",
  "    displayName: Status",
  "  field:",
  "    displayName: Field",
  "views:",
  "  - type: table",
  "    name: All resources",
  "    order:",
  "      - file.name",
  "      - tag",
  "      - status",
  "      - field",
  ""
].join("\n");
fs.writeFileSync(path.join(lib, "Library.base"), base);
console.log("notes", list.length, "hubs", Object.keys(hubs).length, "refused", refused.length);