# Contributor guidance

## Scope and roadmap

The project is documented and planned; it is not an implemented DSM package yet. Follow [the phase roadmap](docs/phases/README.md). Complete a phase's acceptance checks and record the exact NAS, DSM, File Station, package, and CPU versions before relying on its conclusions.

The approved target is DSM 7.4. DSM 7.3 compatibility is a later, unverified concern. Phase 2 is the usable standalone product. File Station integration is Phase 3 and remains conditional.

## Runtime and dependencies

- The only externally installed runtime dependency may be SynoCommunity `mediainfo`.
- First test whether DSM's built-in PHP can safely provide the package request path. Record the executable, SAPI, version, extensions, account visibility, authentication route, and support status. Do not assume an observed PHP version is a public third-party runtime contract.
- If PHP does not meet that gate, bundle a small Go executable built with the standard library for each tested architecture. Do not require Go to be installed on the NAS.
- Build tooling, the Synology Package Toolkit, and CI tooling are build-time concerns, not installed package dependencies.
- Do not add Web Station, Composer, Symfony, a database, or a frontend framework dependency for the initial package.

## Security invariants

- Run the service as its package account. Never add a root daemon, setuid helper, broad share permission, or global DSM/nginx patch to make media readable.
- Authenticate every request, check authorization for every selected file, and do not accept a pathname or username supplied by the browser as proof of access.
- Invoke MediaInfo using its absolute package path, a fixed argument array, and a minimal child environment. Never construct shell commands from filenames, tags, options, or request data.
- Treat media metadata as untrusted. Bound input, process runtime, memory, output, concurrency, temporary storage, and cancellation; clean temporary data on every lifecycle path.
- Prefer proven direct access to the original authorized file. A complete private temporary copy is allowed only as a measured diagnostic prototype, never as an assumed production design for large media.

## DSM and File Station boundaries

- Use published Synology interfaces where they exist. Record whether each integration claim is **documented**, **observed in a legacy package**, or a **hypothesis**.
- `fb_extern` was observed in legacy Synology packages. It is not a documented DSM 7.4 extension contract. Any Phase 3 implementation must be experimental, opt-in, version-scoped, and removable without affecting the standalone app.
- Do not patch, replace, decrypt, redistribute, or inject into File Station or DSM system files. Do not modify the File Station Properties dialog; no supported tab hook has been found.
- Keep File Station launch eligibility client-side only as a usability aid. The backend remains the authorization boundary.

## Documentation and evidence

- For framework, API, SDK, CLI, or library questions, resolve and query current documentation with Context7 first. When it has no authoritative Synology material, use the relevant primary Synology documentation and record that limitation.
- Link claims to primary sources where possible. Preserve research dates and phrase past observations as observations, not current guarantees.
- Exercise behavior on a disposable NAS: authentication and expiry, DSM ACL allow/deny cases, malformed files and hostile tags, path and symlink races, resource limits and cancellation, MediaInfo absence/upgrade, package install/upgrade/uninstall, and File Station enable/disable where applicable.
- Report incomplete validation plainly. Do not add proprietary Synology code, downloaded package contents, credentials, media samples, or session material to this repository.

## Git workflow

- `main` is the only long-lived branch. Do not develop directly on it.
- Start work from an up-to-date `main` and create a short-lived `feature/`, `fix/`, `chore/`, `refactor/`, or `docs/` branch.
- Push the branch and open a pull request targeting `main`. Ensure CI is current and all conversations are resolved before merge.
- Use squash merge and delete the branch afterward.
- Commits on non-`main` branches must be verified signed. The resulting squash commit on `main` is not required to be signed.
