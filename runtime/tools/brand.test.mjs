import test from "node:test";
import assert from "node:assert/strict";
import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";

const brand = path.resolve(path.dirname(fileURLToPath(import.meta.url)), "../../assets/brand");

function attr(file, name) {
  const text = fs.readFileSync(path.join(brand, file), "utf8");
  const match = text.match(new RegExp(`${name}="([^"]+)"`));
  assert.ok(match, `${file} has ${name}`);
  return match[1];
}

test("social preview is 1280 by 640", () => {
  assert.equal(attr("social.svg", "width"), "1280");
  assert.equal(attr("social.svg", "height"), "640");
});

test("the mark is drawn at 16, 64, and 512", () => {
  assert.equal(attr("ow-mark-16.svg", "width"), "16");
  assert.equal(attr("ow-mark.svg", "width"), "64");
  assert.equal(attr("ow-mark-512.svg", "width"), "512");
});

test("three concepts and the default mark are on disk", () => {
  for (const file of ["ow-mark-bars.svg", "ow-mark-ring.svg", "ow-mark-blocks.svg", "ow-mark.svg"]) {
    assert.equal(fs.existsSync(path.join(brand, file)), true, file);
  }
});
