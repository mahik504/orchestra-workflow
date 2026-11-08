# SaaS checklist

Use this when the product has more than one tenant. A single-user script can skip it.

- A cross-tenant test fails when tenant A can read or write tenant B's row. Row-level security, or the equivalent query filter, is what the test calls.
- Zod (or the same schema check) runs at every server action and API boundary. The handler does not trust the client shape.
- One error envelope. The client can tell a bad request from a server failure. Do not invent a second error shape for one route.
- A retried write carries an idempotency key. The second call returns the first result.
- A correlation id is logged on the request and on the error. Tokens, passwords, and session cookies are not logged.
