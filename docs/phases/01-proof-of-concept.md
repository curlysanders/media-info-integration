# Phase 1 — Proof of concept

## Objective

Establish whether a DSM 7.4 package can securely analyze one user-authorized media file with SynoCommunity MediaInfo, and whether package-owned `fb_extern` registration is accepted by File Station without DSM file changes.

## Prerequisites

- A disposable DSM 7.4 NAS with recorded model, CPU architecture, DSM build, File Station version, and test accounts with contrasting ACLs.
- SynoCommunity MediaInfo installed from a recorded version and architecture build.
- Test media that includes a normal file, malformed input, large MKV/MP4, hostile metadata, and access-denied paths.

## Work

### 1. Establish the package request runtime

Inventory the built-in PHP executable, version, SAPI, enabled extensions, package-account access, and a package-owned authenticated request route. Validate that identity is derived from the DSM request environment and that session expiry, logout, CSRF, origin checks, and proxy HTTPS behavior work.

An observed PHP installation is not itself a supported runtime guarantee. Do not use the CLI as an HTTP handler or PHP's built-in development server in production. If PHP cannot meet the route and security requirements without another installed dependency, replace it with a small static Go standard-library executable built for the tested architecture.

### 2. Prove authorization and analysis safety

Expose one authenticated operation that accepts a validated share-relative selection and returns controlled JSON or a controlled error. Prove that the effective requesting DSM user can read the exact input at analysis time, including explicit ACL deny cases. A daemon's ability to read a share is insufficient.

Use `/var/packages/mediainfo/target/bin/mediainfo` with fixed arguments including `--Output=JSON`, an argument array, and a minimal child environment. Reject directories, devices, FIFOs, remote/cloud mounts, URL input, playlist/reference input, and unsafe symlink behavior. Bound file size, runtime, CPU, memory, output, jobs, temporary space, and cancellation.

Preferred evidence is secure access to the original authorized file. A full private temporary copy may be used only to diagnose routing and MediaInfo behavior; measure its space and I/O cost and do not present it as the production solution for large video files.

### 3. Test File Station compatibility separately

On the same disposable target, inspect how DSM 7.4 handles a package-owned `fb_extern` declaration. Determine whether a third-party package is eligible, callback bundle loading and argument shape, selected path semantics, application privilege and lifecycle behavior, and enable/disable/uninstall effects. Register one-file media eligibility only; do not take over double-click preview behavior.

Do not patch File Station, inject browser JavaScript, or alter its Properties dialog. The legacy mechanism is an observed compatibility candidate, not a supported API.

## Deliverables

- A short environment inventory and runtime decision: built-in PHP or bundled Go fallback.
- A reproducible one-file request demonstration with representative MediaInfo JSON.
- Authorization evidence for allow, deny, expiry, tampered path, file replacement/symlink race, malformed media, hostile tags, cancellation, and resource exhaustion.
- A version-scoped `fb_extern` compatibility report stating **documented**, **observed**, or **unavailable**.

## Acceptance checklist

- [ ] The package route authenticates the DSM user and rejects expired or forged requests.
- [ ] The requested user can analyze an allowed file and cannot obtain metadata for a denied file.
- [ ] MediaInfo runs without shell construction or elevated privileges and its output is bounded and safely returned.
- [ ] Temporary data is private and cleaned on success, error, timeout, cancellation, restart, and uninstall.
- [ ] The exact DSM/File Station/MediaInfo/runtime versions and test results are documented.
- [ ] The File Station result is recorded independently; its absence does not block Phase 2.

## Blockers and exclusions

Inability to preserve effective user read authorization blocks Phase 2. An unavailable `fb_extern` route only blocks Phase 3. Do not add a root broker, setuid helper, broad share grant, Web Station dependency, browser companion, DSM patch, or Properties-tab modification to bypass either result.
