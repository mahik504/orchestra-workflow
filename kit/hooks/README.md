# Hooks

The installed hooks live in `kit/git-hooks/` because this repo sets `core.hooksPath` there.

- `commit-msg` removes the Cursor co-author trailer.
- `pre-commit` runs `runtime/tools/hygiene.js` and `runtime/tools/check-versions.js`.

`kit/install-git-hooks.ps1` points the clone at that directory. Do not copy the scripts into a second folder.
