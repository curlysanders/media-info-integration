# P1-004 evidence template — exact-file authorization

Status: **planned live validation**. This report is complete only after the evidence below identifies either a passing documented original-file handoff or a version-scoped unavailable result. Do not represent a package-account filesystem check, pathname precheck, or private copy as user authorization.

## Target and source inventory

| Field | Observed value | Evidence status |
| --- | --- | --- |
| Date/time and operator/method | Pending | Pending |
| DSM model/build and package architecture | Pending recheck | Pending |
| File Station and MediaInfo package/executable versions | Pending recheck | Pending |
| Package revision, SPK digest, UID/GID | Pending | Pending |
| Context7 DSM query/result | Pending | Pending |
| Primary source URLs, publication/version dates | Pending | Pending |

## Candidate decision

| Candidate | Published third-party contract | Identity source | Effective-ACL enforcement and returned object | Decision |
| --- | --- | --- | --- | --- |
| P1-003 custom CGI + `authenticate.cgi` | Identity only; does not establish file authorization | Pending | Pending | Pending |
| `SYNO.FileStation.Download` | Requires an account/password-created `FileStation` session | Pending | Streams content for its own authenticated session | Pending |

State **unavailable** here if no candidate supplies an original-file handle or stream under the authenticated user's effective authorization without an excluded dependency or privilege. Do not fill the live matrix in that case.

## Recreated ACL and package-account matrix

| Selection | `test1` File Station result | `test2` File Station result | Package-account access | Notes |
| --- | --- | --- | --- | --- |
| `shared` fixture | Pending | Pending | Pending | |
| `a-only` fixture | Pending | Pending | Pending | |
| `b-only` fixture | Pending | Pending | Pending | |
| `explicit-deny` fixture | Pending | Pending | Pending | Include otherwise applicable allow |

## Authorization and selection matrix

Record only case labels, status, controlled response shape, and non-sensitive fixture fingerprints. Do not include tokens, cookies, absolute storage paths, filenames returned to a denied caller, or media content.

| Case | User | Expected result | Observed result | Evidence status |
| --- | --- | --- | --- | --- |
| Allowed `shared` selection | `test1` | Original authorized handle/stream | Pending | Pending |
| Allowed `shared` selection | `test2` | Original authorized handle/stream | Pending | Pending |
| `a-only` deny | `test2` | Deny without metadata | Pending | Pending |
| `b-only` deny | `test1` | Deny without metadata | Pending | Pending |
| `explicit-deny` | `test2` | Deny without metadata | Pending | Pending |
| Tampered/traversal/absolute/URL/reference input | Both | Invalid request without metadata | Pending | Pending |
| Directory/device/FIFO/mount/cloud/symlink | Both | Rejected before open | Pending | Pending |
| Replacement/symlink race | Both | No stale authorization | Pending | Pending |
| Rename/deletion/ACL revocation/change during operation | Both | Controlled documented outcome | Pending | Pending |

## Decision

Pending. A pass requires a documented candidate and successful matrix proving that user identity, effective authorization, and one original opened object are bound without privilege elevation. Otherwise state the exact tested target as **unavailable**, identify P1-005 and Phase 2 as blocked, and leave the P1-004 authorization checks unchecked.
