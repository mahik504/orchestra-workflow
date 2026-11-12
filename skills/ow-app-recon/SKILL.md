---
name: ow-app-recon
description: Map an authorized web, Electron, Android, desktop, or iOS app the operator owns, built for a client, or that is a CTF or open source. No decryption and no Frida.
---

# App recon

Authorized targets only. A written `.orchestra/scope.md` names the target and why it is in scope.

## Preconditions

- The target is yours, a client's, a CTF, or open source. Say which.
- `.orchestra/scope.md` exists before any unpack or decompile.
- `orchestra doctor --recon` reports Java, jadx, apktool, and MobSF. It does not install them.

## Procedure

1. Web: save a HAR from a session you are allowed to run. Draft OpenAPI and a state map from that HAR.
2. Electron: list the asar contents and note the bundle layout. Recover sourcemaps only for your own bundle or an open-source one.
3. Android: static analysis with jadx, apktool, and MobSF when those tools are on PATH.
4. Native desktop: load `rea` only when the doctor shows Java on PATH.
5. iOS: record screens and flows on a simulator or device you own. Do not decrypt. Do not attach Frida.

## Output

A scope file, a map of screens or endpoints, and a list of tools that actually ran. Unknowns stay unknown.

## Verify

Every artifact points at the scope file. No third-party sourcemap was reconstructed.

## Skip

Stop if the scope file is missing, the target is someone else's product, or a tool is not installed. Do not install it yourself.
