---
name: ow-scrape
description: Extract public records with a schema, a rate limit, and a provenance log. Stop at login, a paywall, a CAPTCHA, robots, or a 403 or 429.
---

# Scrape

Public data only. No LinkedIn or Instagram. No proxy rotation to dodge a block.

## Preconditions

- A `policy.json` names the domains, a rate of at most one request per second per domain, the user agent, robots compliance, and the cache path.
- The fields you want are a Zod schema or a JSON Schema before the first request.
- The source is an official API, feed, sitemap, or public HTML page.

## Procedure

1. Prefer the official API, feed, sitemap, or JSON-LD.
2. Otherwise fetch and parse.
3. Use Scrapling by default for HTML that needs a real parser. Crawl4AI is the fallback. Playwright is for pages that need JavaScript. Firecrawl only when a key is already in the environment.
4. Write one JSONL provenance row per response: url, time, sha256, extractor.
5. Resume, dedupe, and export JSONL or CSV. Write `SCRAPE_REPORT.md`.

## Output

Rows that match the schema, the provenance log, and the report. Empty cells stay empty.

## Verify

The report names the domains, the rate, and any URL that was skipped.

## Skip

Stop on a login wall, paywall, CAPTCHA, robots disallow, or a 403 or 429 after one backoff. Do not solve the challenge.
