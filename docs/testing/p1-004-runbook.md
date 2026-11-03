# P1-004 — exact-file authorization decision gate

This procedure is for the disposable `DSM-7.4-Test` VM only. It starts from the protected custom-CGI route established by P1-003 and determines whether a standalone package can analyze the **original** file selected by the authenticated DSM user without privilege elevation. It does not authorize MediaInfo execution; P1-005 remains blocked until this gate passes.

## 1. Non-negotiable pass condition

A passing mechanism must, in one bounded request flow:

1. derive the caller from the P1-003 custom-CGI authentication check;
2. apply DSM's effective permission decision for that caller to one selected file at open time; and
3. supply the package only the resulting original-file handle or stream.

The package account must not obtain broad read access to the test share. A pathname permission check followed by an open as the package account is not sufficient: it can disagree with DSM ACLs and permits replacement or symlink races.

Do not use a root broker, setuid helper, global DSM/nginx modification, Web Station dependency, browser companion, copied browser credentials, ACL reimplementation, or an undocumented legacy route. A private complete copy is outside this ticket: it is not a Phase 2 substitute and the user selected no diagnostic-copy fallback.

## 2. Recreate and record the disposable target

Before testing, recreate the removed P1-001 resources and re-record their actual versions:

- DSM model, build, package architecture, File Station version, SynoCommunity MediaInfo package/executable version, package revision/SPK digest, and package UID/GID;
- ordinary non-administrator users `test1` and `test2`, including group membership;
- share `mediainfo-p1-test`, with `shared`, `a-only`, `b-only`, and `explicit-deny` folders; and
- the P1-001 fixture descriptions and checksums, plus a distinct-content replacement fixture for race checks.

Use the P1-001 matrix: both users can read `shared`; only `test1` can read `a-only`; only `test2` can read `b-only`; `test2` has an explicit deny on `explicit-deny` despite an otherwise applicable allow. Re-verify those results in independent File Station sessions. Add a controlled deny where the package account could otherwise read, but do not grant the account the whole share.

## 3. Documentation and candidate gate

On the day of testing, resolve and query Context7 first for DSM 7 third-party package/file-authorization documentation. Record its result, including a no-match result. Then consult dated primary Synology material.

The [DSM Developer Guide](https://global.download.synology.com/download/Document/Software/DeveloperGuide/Os/DSM/All/enu/DSM_Developer_Guide_7_enu.pdf) documents invoking `authenticate.cgi` from a package custom CGI to identify an authenticated caller. It does not, by itself, document per-user filesystem authorization or an original-file descriptor handoff. Synology's [File Station ACL and Permission Inspector guidance](https://kb.synology.com/en-eu/DSM/tutorial/How_to_manage_ACL_settings_on_your_NAS) documents configuring and viewing permissions, not a third-party package API that opens a file as the selected user.

The official [File Station API Guide](https://global.download.synology.com/download/Document/Software/DeveloperGuide/Package/FileStation/All/enu/Synology_File_Station_API_Guide.pdf) documents `SYNO.FileStation.Download` returning file content, but requires a separately logged-in `SYNO.API.Auth` session with `session=FileStation`; its documented login request requires an account and password. It does not document deriving that session from a package custom-CGI request or forwarding the custom-CGI request cookie to obtain it. Therefore it is not a passing candidate unless Synology publishes that missing handoff contract.

Make a candidate table for every published DSM interface found. For each candidate, record its source URL/version/date, third-party package support, request identity source, ACL enforcement point, exact object returned, execution account, installed dependencies, and whether it preserves the original object without a second pathname open. Reject a candidate immediately if any pass-condition item is missing.

Only a documented candidate may be installed or exercised. Reading documentation and inspecting package-owned metadata are safe research steps; do not test an observed legacy mechanism as a fallback. If no candidate survives this table, record **unavailable** for the exact tested target, skip the remaining live matrix, and block P1-005 and Phase 2.

## 4. Bounded selection and exact-object proof

Implement a test-only `authorize-selection` operation only after a candidate passes the documentation gate. It must remain POST-only and retain P1-003's authentication, CSRF/token, origin, request-size, and no-store response protections.

Its request represents exactly one share-relative file selection. Bound all field and total lengths; use UTF-8 text; require the fixed test share name and a non-empty relative path. Reject absolute paths, empty or `.`/`..` components, separators that create an alternate path syntax, URLs, query/fragment input, playlists/references, and a NUL byte. Do not reveal the normalized path, filename, size, type, existence, or content to a denied caller.

The candidate must authorize and obtain one immutable original-file handle or stream before any result can be returned. It must reject directories, devices, FIFOs, remote/cloud mounts, and symlinks unless the published interface proves an equivalent safe object binding. The package must neither reopen the selected pathname nor convert a permission result into package-account access. Keep this operation free of MediaInfo execution and file metadata responses.

## 5. Live matrix and decision

For every permitted candidate, run each case in separate `test1` and `test2` DSM sessions. Record only a case label, status, controlled response shape, source/route version, and a fixture fingerprint observed out of band; never store cookies, tokens, absolute storage paths, or media content.

| Case | Expected result |
| --- | --- |
| `shared` allowed selection for each user | Authorized original handle/stream; no package share grant |
| `a-only`, `b-only`, and `explicit-deny` opposite-user selections | Controlled deny with no file metadata |
| Malformed, traversal, absolute, URL/reference, and ambiguous selections | Controlled invalid-request response with no file metadata |
| Directory, device, FIFO, mount, cloud/remote target, and symlink | Controlled rejection before content is opened |
| Replacement or symlink inserted between validation and open | Reject/retry; never consume the replacement through a stale authorization |
| Rename, deletion, and ACL revocation before open | Controlled deny/not-found response; no content or metadata leak |
| Rename, deletion, ACL revocation, or write after the candidate opens the object | Record exact DSM behavior and prove the consumed object remains the authorized handle/stream, or reject the candidate |

For race cases, alternate known distinct fixture contents and retain their fingerprints only in the sanitized evidence. Run enough iterations to demonstrate the ordering control used; do not claim a probabilistic race test proves an atomic guarantee.

P1-004 passes only when the candidate meets every pass condition and every required matrix row has a controlled outcome. Update the ticket checkboxes and evidence only after that result. Otherwise publish a version-scoped **unavailable** result, leave the authorization acceptance checks unchecked, and record P1-005 plus Phase 2 as blocked. This negative decision completes P1-004 research; it does not authorize a workaround.

## 6. Cleanup

Remove the P1-004 package build, test share, users, fixtures, and temporary registrations using the P1-001 reset procedure. Preserve the sanitized evidence report and no credentials, session material, DSM code, downloaded package contents, or media samples.
