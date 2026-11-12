---
name: ow-browserize
description: Turn an authorized browser session you already completed into a Playwright script with fixtures and assertions. Never store a password or a cookie jar.
---

# Browserize

The session already happened in an authorized browser. This skill writes the test. It does not log in for you.

## Preconditions

- The flow is on an app you own, a client's, a CTF, or open source.
- You can describe the clicks and the expected text without pasting a secret.

## Procedure

1. Write the steps as a Playwright test in the product repo.
2. Use fixtures for data. Assert on visible text, URL, and network status.
3. Keep credentials out of the file. Read them from the environment at run time if a later human run needs them, and do not commit the values.

## Output

A test file that fails when the flow breaks.

## Verify

Run the test against the local app. Paste the exit code.

## Skip

Do not record a session on a site you are not allowed to automate. Do not store cookies in the repo.
