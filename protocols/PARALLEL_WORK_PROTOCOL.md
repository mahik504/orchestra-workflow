# Parallel work

One conductor. The chat that writes the plan. Other hosts are lanes.

## Lane file

`.orchestra/lanes.md` in the product repo:

| lane | host | branch | worktree | owns | acceptance | status |
| --- | --- | --- | --- | --- | --- | --- |
| cursor | Cursor | main | this repo | shared files | the plan | open |

## Rules

1. One git worktree per extra host: `git worktree add ../repo-ag -b lane/ag/task`.
2. A lane edits only the paths in its capsule `owned_files`. Two capsules that claim the same path conflict, and the conductor resolves it before either lane writes.
3. Shared files belong to the conductor lane. That includes package files and DESIGN.md.
4. Only the conductor writes the vault.
5. The human merges and pushes. No co-author trailer.
