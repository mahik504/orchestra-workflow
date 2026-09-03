# Resource Curation Policy

## 1. Orchestra Active Working Set
This is the core execution context. It must remain small, highly reliable, and directly aligned with verified Orchestra capabilities. Items in this set are evaluated for loading *every* time a task is classified.

## 2. Curated Optional Catalog
High-value capabilities (e.g., `r3f-threejs`, 3D shaders, obscure database schemas) live here. They are explicitly NOT loaded by default. They are retrieved *only* when a task manifest explicitly requests that capability. 

## 3. The Quarantine Zone
The 1,598-skill directory (`skills_library`) is strictly quarantined. It serves as an archive/sandbox for raw, untested AI materials. 
**Rule:** The `skills_library` must never be automatically exposed to the active routing context.

## Deletion Criteria
A resource is deleted from Active/Curated if it is:
- Broken (e.g., `github-mcp-server` Docker failure).
- Duplicative (e.g., three different Next.js prompt wrappers).
- Generic AI Slop (e.g., "Write clean code").
