# Phase 1 tickets — Proof of concept

These tickets implement the [Phase 1 roadmap](01-proof-of-concept.md). P1-001 is in progress on the disposable `DSM-7.4-Test` VM; later tickets remain **Planned**. Results are scoped to that Arc Loader VM, not physical Synology hardware. DSM 7.4 is the target. DSM 7.3 compatibility and the usable standalone product belong to later work.

## Execution and evidence

Run NAS experiments only on the disposable target established by P1-001. Run the service as its package account. The only externally installed runtime dependency may be SynoCommunity MediaInfo; do not introduce Web Station, Composer, Symfony, a database, or a frontend framework. Do not add a root broker, setuid helper, broad share grant, or global DSM/nginx changes to bypass a failed gate.

Each ticket must record dated, reproducible steps, expected versus observed results, exact tested versions, and unresolved limitations. Use a sanitized evidence section or linked repository report when executing the ticket; do not add credentials, session material, proprietary code, downloaded package contents, or media samples. Describe fixtures sufficiently to reproduce them without committing them. Leave acceptance boxes unchecked until their evidence exists.

For API, runtime, or tooling documentation, consult Context7 first. If authoritative Synology material is unavailable there, record that limitation and use primary Synology sources. Label integration claims as **documented**, **observed**, or **hypothesis**; past observations are not current guarantees. Starting references are the [feasibility report](../feasibility.md) and [legacy-package investigation](../package-integration-research.md).

| Ticket | Outcome | Dependencies | Status |
| --- | --- | --- | --- |
| [P1-001](#p1-001--establish-the-disposable-test-environment) | Environment inventory and reproducible fixtures | — | In progress |
| [P1-002](#p1-002--select-and-prove-the-package-request-runtime) | PHP assessment and runtime decision | P1-001 | Planned |
| [P1-003](#p1-003--prove-authenticated-package-routing) | Authenticated and protected request route | P1-002 | Planned |
| [P1-004](#p1-004--prove-user-authorization-for-the-exact-opened-file) | Effective user authorization and safe file access | P1-003 | Planned |
| [P1-005](#p1-005--demonstrate-controlled-one-file-mediainfo-analysis) | Reproducible one-file JSON demonstration | P1-004 | Planned |
| [P1-006](#p1-006--validate-containment-resource-limits-and-cleanup) | Measured limits, cancellation, and cleanup | P1-005 | Planned |
| [P1-007](#p1-007--investigate-file-station-compatibility) | Independent, version-scoped compatibility result | P1-001 | Planned |
| [P1-008](#p1-008--consolidate-evidence-and-record-readiness-decisions) | Separate Phase 2 and Phase 3 decisions | P1-001–P1-007 results | Planned |

## P1-001 — Establish the disposable test environment

**Objective:** Make every subsequent observation attributable to a known target and reproducible test setup.

**Dependencies:** None.

**Work:** Inventory the disposable DSM 7.4 test VM and installed MediaInfo package. Prepare contrasting user ACLs and local test fixtures without granting the service broad share access. Keep host observations distinct from guest-verified facts.

**Acceptance criteria:**

- [ ] Record presented NAS model, guest CPU model and package architecture, exact DSM build, File Station version, MediaInfo package version and architecture build, and MediaInfo executable version. Record VM and loader details so results remain version-scoped.
- [x] Record the prototype package version/revision and package-account identity, or explicitly mark them pending until the prototype exists and update them before interpreting its results.
- [ ] Prepare at least two DSM users with contrasting read permissions, including explicit ACL deny cases, and record the intended allow/deny matrix.
- [ ] Prepare short MKV and MP4 videos with identifiable video and audio codecs, plus malformed input, hostile metadata, and access-denied paths; record track details, sizes, checksums, and reproducible fixture descriptions.
- [ ] Record safe reset and cleanup steps for test package registrations, files, and accounts on the disposable target.

**Required evidence:** Dated environment inventory, fixture descriptions, ACL matrix, and reset procedure. The [P1-001 runbook](../testing/p1-001-runbook.md) and [dated evidence](../testing/p1-001-evidence-2026-10-26.md) track this work. Later tickets reference the inventory and record any version changes.

## P1-002 — Select and prove the package request runtime

**Objective:** Establish a package-owned request path without another installed runtime dependency.

**Dependencies:** P1-001.

**Work:** Test DSM's built-in PHP first. Assess whether it can host the package request route and support the security checks in P1-003. If it cannot meet that gate, prove a small bundled Go executable using only the standard library for the tested architecture. Revisit the decision if P1-003 exposes a runtime limitation.

**Acceptance criteria:**

- [ ] Record PHP executable, version, SAPI, enabled extensions, visibility/execution under the package account, candidate authentication route, and documented support status.
- [ ] Demonstrate a package-owned request route under the package account without global DSM/nginx edits or another installed runtime dependency.
- [ ] Distinguish an observed PHP installation from a documented third-party runtime contract; do not use CLI PHP as an HTTP handler or the PHP development server as the production route.
- [ ] Record the PHP gate outcome and rationale. If it fails, build and demonstrate a bundled static Go standard-library executable for the tested architecture, recording compiler/build details without requiring Go on the NAS.
- [ ] Record selected runtime version, prototype package revision, routing mechanism, and unresolved security assumptions for P1-003.

**Required evidence:** Runtime inventory, route demonstration, primary-source support assessment, and explicit PHP-or-Go decision. No untested architecture is claimed supported.

## P1-003 — Prove authenticated package routing

**Objective:** Derive the requesting DSM identity from a verified request environment and reject forged or expired requests.

**Dependencies:** P1-002.

**Work:** Exercise the selected package route with authenticated POST requests, application access policy, CSRF protection, origin checks, and bounded input. Validate the route through direct HTTPS and the intended reverse proxy configuration.

**Acceptance criteria:**

- [ ] Demonstrate identity verification from the DSM request environment; browser-supplied usernames or public headers cannot impersonate another user.
- [ ] Verify valid sessions and rejection of unauthenticated, forged, expired, and logged-out sessions on every request.
- [ ] Verify application access allow/deny behavior independently of menu visibility.
- [ ] Reject missing/invalid CSRF tokens, disallowed origins, unsupported methods, and oversized or malformed request bodies.
- [ ] Verify HTTPS/proxy handling without trusting unvalidated forwarded headers as identity or security proof.
- [ ] Keep credentials/session material out of URLs, logs, persistent storage, and child arguments; return controlled errors and session-scoped responses with `Cache-Control: no-store`.

**Required evidence:** Sanitized request/result matrix for both users and route configurations, including expiry/logout and forgery cases. If the runtime cannot satisfy these checks, update P1-002 before continuing.

## P1-004 — Prove user authorization for the exact opened file

**Objective:** Bind authenticated user read permission to the exact input analyzed, without privilege elevation.

**Dependencies:** P1-003.

**Work:** Accept a validated share-relative selection and investigate secure direct access to the original file. Prove effective DSM ACL behavior at analysis time; a service account's read access or a pathname permission check alone is insufficient.

**Acceptance criteria:**

- [ ] Both users can access allowed files and receive no file metadata for denied files, including explicit ACL deny cases where the service account could otherwise read the file.
- [ ] Document how the authenticated identity, effective read authorization, and exact opened file are bound without a root broker, setuid helper, or broad share permission.
- [ ] Reject tampered selections, traversal, directories, devices, FIFOs, remote/cloud mounts, URL input, playlist/reference input, and unsafe symlink behavior.
- [ ] Exercise replacement/symlink races, rename, deletion, access revocation, and files changing during analysis; record controlled outcomes and remaining limitations.
- [ ] Prove access to the original authorized file. If a complete private copy is used diagnostically, label it diagnostic-only and measure size, additional space, I/O, and latency on the prepared fixtures; do not count it as the Phase 2 direct-access gate passing.
- [ ] Record a blocked Phase 2 gate if effective user authorization cannot be preserved; do not substitute the daemon's access as proof.

**Required evidence:** Access-design explanation, allow/deny matrix, race-test procedure/results, and any diagnostic-copy measurements. A negative finding must remain explicit rather than being marked as successful authorization.

## P1-005 — Demonstrate controlled one-file MediaInfo analysis

**Objective:** Produce representative JSON through the authenticated, authorized path with controlled failure behavior.

**Dependencies:** P1-004. A failed authorization gate blocks an end-to-end success claim; diagnostic parser experiments must be identified separately.

**Work:** Demonstrate one authenticated operation taking a validated share-relative selection and returning controlled JSON or an error. Invoke the installed MediaInfo binary under the package account and treat its output as untrusted.

**Acceptance criteria:**

- [ ] Invoke `/var/packages/mediainfo/target/bin/mediainfo` using a fixed argument array including `--Output=JSON`, one verified authorized input, and a minimal child environment without request secrets.
- [ ] Verify executable permissions and library resolution under the package account; record behavior when MediaInfo is absent, unavailable, or upgraded.
- [ ] Demonstrate no shell construction, user-supplied CLI options, or elevated privileges, including filenames containing spaces, shell characters, and option-like prefixes.
- [ ] Validate JSON shape and bound stdout/stderr; malformed media, invalid JSON, nonzero exits, and oversized output yield controlled errors without sensitive paths or diagnostics leaking.
- [ ] Return hostile metadata safely as data, replace local paths with the authorized display name, and scope results to the requesting session. Any demonstration display renders strings as text.
- [ ] Record a reproducible normal-file request and sanitized representative JSON, plus malformed-input and hostile-tag results. Record the initial enforced limits for P1-006 to exercise.

**Required evidence:** Invocation/environment description, sanitized request/output examples, dependency failure results, and error/output validation matrix.

## P1-006 — Validate containment, resource limits, and cleanup

**Objective:** Demonstrate bounded parsing and complete cleanup under success and failure conditions.

**Dependencies:** P1-005.

**Work:** Establish concrete enforced limits from measurements on the target. Exercise parser isolation, resource exhaustion, cancellation, and lifecycle cleanup, including any diagnostic temporary copies.

**Acceptance criteria:**

- [ ] Document and exercise numeric limits for input/file size, wall time, CPU, memory, stdout/stderr, concurrency, and temporary storage, with boundary and over-limit results.
- [ ] Verify cancellation and timeout terminate the analysis and child processes, release job capacity, and prevent abandoned processing.
- [ ] Verify package-account execution, restricted writable storage, per-job filesystem isolation, separation from other users' inputs, and network restrictions; test media references to other files or network resources. Record unavailable containment facilities as unresolved gates.
- [ ] Test malformed media, hostile tags, concurrent requests, memory/resource exhaustion, and insufficient temporary space; the request path remains controlled and recovers capacity.
- [ ] Keep temporary data private and inaccessible through the web route; verify cleanup after success, error, timeout, cancellation, restart, and uninstall. If no temporary media is created, demonstrate that and still account for other job artifacts.
- [ ] Measure MKV/MP4 latency and storage use on the prepared short clips, distinguish direct analysis from diagnostic-copy costs, and state that large-media performance remains unmeasured.

**Required evidence:** Limit values and enforcement mechanisms, measurement method/results, cancellation/process checks, isolation results, and cleanup matrix with unresolved gaps plainly stated.

## P1-007 — Investigate File Station compatibility

**Objective:** Determine whether a third-party package-owned `fb_extern` declaration works on the tested DSM 7.4 target, independently of the standalone analysis path.

**Dependencies:** P1-001. May run independently of P1-002–P1-006.

**Work:** Test an original package-owned registration and callback bundle on the disposable NAS. Use the legacy research as an investigation lead, without copying proprietary implementations or treating it as a supported extension contract.

**Acceptance criteria:**

- [ ] Determine third-party eligibility, callback bundle loading, callback argument/record shape, and selected path semantics, including share-relative versus physical/virtual paths.
- [ ] Limit client-side eligibility to one supported media file; record directory and multiple-selection behavior, and preserve existing double-click preview behavior.
- [ ] Exercise application privilege allow/deny, package enable/disable, restart, upgrade, uninstall, and removal of the registration; record File Station refresh/session effects and any stale actions.
- [ ] Keep the backend as the authorization boundary; menu eligibility or callback paths do not grant access. Compatibility can be tested with a launch-only prototype before analysis is available.
- [ ] Perform no DSM/File Station patching, browser JavaScript injection, decryption, or Properties-dialog modification; do not commit proprietary code or downloaded package contents.
- [ ] Publish a version-scoped result of **documented**, **observed**, or **unavailable**, with primary sources where applicable and hypotheses/limitations separated from observations. An unavailable result completes this investigation when evidenced.

**Required evidence:** Exact NAS/DSM/File Station/prototype versions, sanitized callback shape and path observations, eligibility/lifecycle matrix, removal verification, and compatibility conclusion. An unavailable route blocks Phase 3 only.

## P1-008 — Consolidate evidence and record readiness decisions

**Objective:** Make the Phase 1 result reviewable and decide each subsequent phase independently.

**Dependencies:** Results from P1-001–P1-007, including negative findings or explicit blocked/unperformed checks. Do not require every investigation to succeed before documenting the outcome.

**Work:** Consolidate the inventory, runtime decision, one-file demonstration, security/containment results, and compatibility report. Map evidence to the Phase 1 acceptance checklist and record separate readiness decisions.

**Acceptance criteria:**

- [ ] Link dated evidence for each ticket, exact tested NAS/CPU/DSM/File Station/MediaInfo/runtime/prototype versions, and reproducible steps with expected versus observed results.
- [ ] Map every Phase 1 acceptance check to evidence; leave unverified or failed checks unchecked and explain blockers and validation gaps.
- [ ] Record Phase 2 as ready only when safe routing, effective user authorization for the exact input, controlled analysis, resource bounds, and cleanup are demonstrated. Explicitly block Phase 2 if user authorization cannot be preserved; diagnostic full copies do not satisfy that gate.
- [ ] Record the independent Phase 3 compatibility decision and its tested version scope. File Station unavailability does not block an otherwise ready Phase 2.
- [ ] Update Phase 1/roadmap statuses to match actual evidence, without claiming an implemented standalone product or untested architecture/DSM compatibility.
- [ ] Verify repository evidence contains no credentials, sessions, proprietary package contents, or media samples, and that all report links resolve.

**Required evidence:** Consolidated Phase 1 report with acceptance-to-evidence mapping, separate Phase 2/Phase 3 readiness decisions, and explicit remaining blockers.

## Phase 1 acceptance coverage

| Roadmap acceptance check | Primary tickets |
| --- | --- |
| Authenticate the DSM user; reject expired or forged requests | P1-002, P1-003 |
| Analyze allowed files without disclosing denied-file metadata | P1-004, P1-005 |
| Run MediaInfo without shells/elevation; bound and safely return output | P1-005, P1-006 |
| Private temporary data and cleanup on every lifecycle path | P1-006 |
| Exact versions and test results documented | P1-001, all execution evidence, P1-008 |
| Independent File Station result; absence does not block Phase 2 | P1-007, P1-008 |
