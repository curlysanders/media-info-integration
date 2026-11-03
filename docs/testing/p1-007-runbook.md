# P1-007 — File Station `fb_extern` compatibility probe

This runbook tests the package-owned, undocumented File Station integration candidate on the disposable DSM target only. It is an **observed compatibility experiment**, not a supported Synology extension contract and not an authorization mechanism.

The `0.1.0-0022` probe contributes a minimal `fb_extern` declaration and callback bundle. It accepts exactly one regular file whose extension is one of `m4v`, `mkv`, `mov`, `mp4`, `mpeg`, `mpg`, or `webm`. It recognizes an array, Ext collection, `{records}`, `{recs}`, or a single record callback shape. On activation it opens a package-owned window and displays the File Station record as a client-side diagnostic. It does not submit the record to the CGI, open a file, invoke MediaInfo, relay a session, or persist any selection.

## Preconditions

- Disposable target: DSM `7.4.1-90080 Update 1`, `v1000`, File Station `1.4.4-2221`.
- A package build that preserves the P1-003 local routing policy rather than using the committed example policy.
- One harmless generated media fixture in `mediainfo-p1-test`; do not upload user media.
- `test1` and `test2` have independent File Station sessions and the ACL matrix from P1-004.

## Procedure

1. Record the SPK SHA-256 and package version. Upgrade the existing disposable `mediainfo-poc` package through Package Center.
2. Open a fresh File Station session. If necessary, sign out/in or reload File Station so it reads package registrations; do not patch or reload DSM/File Station files.
3. Right-click one supported media fixture and record whether **View MediaInfo (experimental)** appears. Activate it and verify that the menu closes and a package-owned window displays the fixture name and File Station reference.
4. Repeat with a text file, directory, and multiple selection. The action must be absent or disabled in every ineligible case.
5. Disable and re-enable the package, then repeat the eligible case. Upgrade the same package once and verify that File Station remains functional.
6. Do not test analysis from this callback yet. P1-004X must separately prove any File Station session relay and streaming behavior before a selected record is accepted by the backend.

## Record

Record DSM, File Station, package, and SPK versions; the menu and window outcome for each case; callback-record shape only after removing user paths; refresh/session behavior; and enable/upgrade effects. Do not record cookies, tokens, passwords, downloaded DSM code, or user media.

## Decision

Mark the mechanism **observed** only if a third-party package registration works without DSM/File Station modifications and removal leaves File Station functional. It remains experimental and version-scoped. An unavailable result completes P1-007 and blocks Phase 3 only.
