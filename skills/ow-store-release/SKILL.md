---
stub: true
stub_reason: Short pointer. See docs/audit/STUB-TRIAGE.md.
name: ow-store-release
description: Checklist for shipping a mobile or desktop app to a store. Use for Play, the App Store, or a Microsoft Store package. Do not claim the app is listed.
---

# Store release

| Store | Check |
| --- | --- |
| Play | Package id, version code, signing, privacy policy, screenshots |
| App Store | Bundle id, version, privacy labels, screenshots, review notes |
| Microsoft Store | MSIX identity, version, signing, age rating, screenshots |

Expo and EAS stay the default path. Native iOS from Windows is a separate pack. Do not say the build is in the store until the user shows the listing.
