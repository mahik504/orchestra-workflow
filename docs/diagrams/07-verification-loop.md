# Verification loop

```mermaid
flowchart TD
  impl[Implement]
  launch[Launch the real app]
  exercise[Exercise the flow]
  stills[Viewports + contrast number]
  sos[Designer SOS]
  hall[Hallmark after stills]
  ledger[Verification ledger]
  impl --> launch --> exercise --> stills --> sos --> hall --> ledger
```

Playwright verifies **our** app. agent-browser is optional for **foreign** sites. Unlazy is methodology, not a second completion authority. No DONE/FIXED/VERIFIED/PASSED/SHIPPED without evidence in the same message.
