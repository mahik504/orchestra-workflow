#!/usr/bin/env node
/**
 * Canonical Layer A capture for Orchestra.
 * Resolves the Playwright Node library from runtime/package.json (not Playwright
 * Test CLI, not @playwright/mcp). Chromium only. Pass/fail is owned by the Go harness.
 */
import { createRequire } from "node:module";
import { fileURLToPath } from "node:url";
import fs from "node:fs";
import path from "node:path";

const here = path.dirname(fileURLToPath(import.meta.url));
const runtimeRoot = path.resolve(here, "..");
const require = createRequire(path.join(runtimeRoot, "package.json"));

function fail(msg, missing = false) {
  const payload = {
    owner: "harness-playwright",
    passed: false,
    environment_error: msg,
    playwright_missing: missing,
    browser_closed: true,
  };
  process.stdout.write(JSON.stringify(payload));
  process.exit(0);
}

let chromium;
try {
  ({ chromium } = require("playwright"));
} catch (err) {
  fail(
    "playwright library missing under runtime/node_modules; run npm install in runtime/",
    true
  );
}

const specPath = process.argv[2];
if (!specPath) {
  fail("layer-a.mjs requires a JSON spec path");
}

let spec;
try {
  spec = JSON.parse(fs.readFileSync(specPath, "utf8"));
} catch (err) {
  fail("invalid Layer A spec: " + err.message);
}

const baseURL = (spec.base_url || "").trim();
if (!baseURL) {
  fail("no running application URL; Layer A cannot invent evidence");
}

const outputDir = spec.output_dir || path.join(".", ".orchestra", "qa");
fs.mkdirSync(outputDir, { recursive: true });
const viewports = Array.isArray(spec.viewports) && spec.viewports.length
  ? spec.viewports
  : [
      { name: "desktop", width: 1440, height: 900 },
      { name: "mobile", width: 390, height: 844 },
    ];
const routes = Array.isArray(spec.routes) && spec.routes.length ? spec.routes : ["/"];

function pngMagic(buf) {
  return (
    buf.length >= 8 &&
    buf[0] === 0x89 &&
    buf[1] === 0x50 &&
    buf[2] === 0x4e &&
    buf[3] === 0x47 &&
    buf[4] === 0x0d &&
    buf[5] === 0x0a &&
    buf[6] === 0x1a &&
    buf[7] === 0x0a
  );
}

function joinURL(base, route) {
  if (!route || route === "/") return base;
  try {
    return new URL(route, base.endsWith("/") ? base : base + "/").toString();
  } catch {
    return base.replace(/\/$/, "") + (route.startsWith("/") ? route : "/" + route);
  }
}

async function settle(page) {
  await page.waitForLoadState("domcontentloaded");
  await page.evaluate(
    () =>
      new Promise((resolve) => {
        requestAnimationFrame(() => requestAnimationFrame(resolve));
      })
  );
}

const t0 = Date.now();
let browser;
let browserClosed = false;
let startupMS = 0;
let screenshotMS = 0;
const result = {
  owner: "harness-playwright",
  passed: true,
  captures: [],
  browser_name: "chromium",
  browser_version: "",
  browser_closed: false,
  browser_pid: 0,
  chromium_path: "",
  startup_ms: 0,
  screenshot_ms: 0,
};

try {
  const chromiumPath = chromium.executablePath();
  if (!chromiumPath || !fs.existsSync(chromiumPath)) {
    fail("Chromium binary missing; run npx playwright install chromium from runtime/", true);
  }
  result.chromium_path = chromiumPath;

  browser = await chromium.launch({
    headless: true,
    handleSIGINT: true,
    handleSIGTERM: true,
    handleSIGHUP: true,
  });
  startupMS = Date.now() - t0;
  result.startup_ms = startupMS;
  result.browser_name = browser.browserType().name();
  result.browser_version = browser.version();
  const proc = typeof browser.process === "function" ? browser.process() : browser.process;
  result.browser_pid = proc && proc.pid ? proc.pid : 0;

  const shotStart = Date.now();
  for (const route of routes) {
    const url = joinURL(baseURL, route);
    for (const vp of viewports) {
      const context = await browser.newContext({
        viewport: { width: vp.width, height: vp.height },
      });
      const page = await context.newPage();
      const consoleErrors = [];
      const networkFailures = [];
      page.on("console", (msg) => {
        if (msg.type() === "error") consoleErrors.push(msg.text());
      });
      page.on("pageerror", (err) => consoleErrors.push(String(err)));
      page.on("response", (res) => {
        if (res.status() >= 400) networkFailures.push(`${res.status()} ${res.url()}`);
      });

      await page.goto(url, { waitUntil: "load", timeout: 20000 });
      await settle(page);

      for (const sel of spec.required_selectors || []) {
        const loc = page.locator(sel);
        if ((await loc.count()) === 0) {
          throw new Error("required selector missing: " + sel);
        }
      }

      for (const step of spec.interactions || []) {
        if (step.kind === "click" && step.selector) {
          await page.locator(step.selector).first().click({ timeout: 5000 });
        } else if (step.kind === "type" && step.selector) {
          await page.locator(step.selector).first().fill(step.value || "", { timeout: 5000 });
        } else if (step.kind === "key" && step.value) {
          await page.keyboard.press(step.value);
        } else if (step.kind === "navigate" && step.value) {
          await page.goto(joinURL(baseURL, step.value), { waitUntil: "load", timeout: 20000 });
          await settle(page);
        }
      }

      const buttons = page.locator("button, [role=button]");
      if ((await buttons.count()) > 0) {
        try {
          await buttons.first().click({ timeout: 2000 });
        } catch {
          /* optional */
        }
      }
      const inputs = page.locator("input:not([type=hidden]), textarea");
      if ((await inputs.count()) > 0) {
        try {
          const first = inputs.first();
          const cur = await first.inputValue().catch(() => "");
          if (!cur) await first.fill("orchestra-layer-a", { timeout: 2000 });
        } catch {
          /* optional */
        }
      }

      const metrics = await page.evaluate(() => {
        const doc = document.documentElement;
        const hero = document.querySelector("[data-qa=hero], h1, header");
        const labels = Array.from(document.querySelectorAll("input, textarea, select")).map((el) => {
          const id = el.getAttribute("id");
          const hasLabel = id && document.querySelector(`label[for="${id}"]`);
          return hasLabel || el.getAttribute("aria-label") ? "ok" : "unlabeled";
        });
        const cs = hero ? getComputedStyle(hero) : getComputedStyle(document.body);
        return {
          scrollWidth: doc.scrollWidth,
          innerWidth: window.innerWidth,
          overflow: doc.scrollWidth > window.innerWidth + 1,
          font: cs.fontFamily,
          size: cs.fontSize,
          visible: !!(hero ? hero.offsetParent !== null || hero.getClientRects().length : true),
          a11y: labels,
        };
      });

      const shotPath = path.resolve(outputDir, `${vp.name}.png`);
      await page.screenshot({ path: shotPath, fullPage: true });
      const buf = fs.readFileSync(shotPath);
      const capture = {
        viewport: vp,
        route,
        screenshot: shotPath,
        overflow: metrics.overflow,
        scroll_width: metrics.scrollWidth,
        console_errors: consoleErrors,
        network_failures: networkFailures,
        visible: metrics.visible,
        computed: { fontFamily: metrics.font, fontSize: metrics.size },
        passed: true,
        failures: [],
        a11y: metrics.a11y,
        byte_size: buf.length,
        png_magic: pngMagic(buf),
        url,
      };
      if (metrics.overflow) {
        capture.passed = false;
        capture.failures.push("overflow");
        result.passed = false;
      }
      if (consoleErrors.length) {
        capture.passed = false;
        capture.failures.push("console");
        result.passed = false;
      }
      if (!capture.png_magic) {
        capture.passed = false;
        capture.failures.push("screenshot");
        result.passed = false;
      }
      result.captures.push(capture);
      await context.close();
    }
  }
  screenshotMS = Date.now() - shotStart;
  result.screenshot_ms = screenshotMS;
} catch (err) {
  result.passed = false;
  result.environment_error = String(err && err.message ? err.message : err);
} finally {
  if (browser) {
    try {
      await browser.close();
      browserClosed = true;
    } catch {
      browserClosed = false;
    }
  } else {
    browserClosed = true;
  }
  result.browser_closed = browserClosed;
}

process.stdout.write(JSON.stringify(result));
