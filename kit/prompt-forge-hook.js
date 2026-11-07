// Reminds Claude Code to run ow-forge. Does not rewrite the prompt.
const out = {
  hookSpecificOutput: {
    hookEventName: "UserPromptSubmit",
    additionalContext: "Orchestra step 0: write a short brief (intent, kind, route, acceptance, refs, constraints). Before any write, print Load: and Skip: from the routing table. A named skill wins. A missing tool stops and is not installed. Ask only if a product decision is ambiguous."
  }
};
process.stdout.write(JSON.stringify(out));
