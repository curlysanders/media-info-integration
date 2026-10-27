# P1-003 evidence template — protected package routing

Status: **planned live validation**. Do not mark complete until every required observation below is populated on the disposable VM.

## Target and artifact inventory

| Field | Observed value | Evidence status |
| --- | --- | --- |
| Date/time and operator/method | Pending | Pending |
| DSM build, model, platform, File Station version | Pending recheck | Pending |
| Package revision, account UID/GID, SPK SHA-256 | Pending | Pending |
| Direct HTTPS origin and certificate result | Pending | Pending |
| DSM Reverse Proxy source/destination and certificate result | Pending | Pending |
| Policy users and origins (non-secret values only) | Pending | Pending |
| Primary source date and Context7 result | Pending | Pending |

## Sanitized result matrix

Record no credentials, cookies, raw tokens, raw headers, URLs containing session material, proprietary DSM output, or media paths.

| Route | Session/user condition | Request mutation | Expected | Observed status/header/body shape | Pass |
| --- | --- | --- | --- | --- | --- |
| Direct HTTPS | `test1` valid | Valid POST and configured origin | 200 | Pending | [ ] |
| Direct HTTPS | `test2` valid | Valid POST and configured origin | 403 | Pending | [ ] |
| Direct HTTPS | No/forged/logout/expired | Valid-shaped POST | 401 | Pending | [ ] |
| Direct HTTPS | `test1` valid | Missing/invalid DSM token | 401 | Pending | [ ] |
| Direct HTTPS | `test1` valid | Forged username/forwarded headers/query | Identity unchanged | Pending | [ ] |
| Direct HTTPS | `test1` valid | Missing/disallowed Origin | 403 | Pending | [ ] |
| Direct HTTPS | `test1` valid | Unsupported method | 405 | Pending | [ ] |
| Direct HTTPS | `test1` valid | Invalid content/body/length | 400 | Pending | [ ] |
| Direct HTTPS | `test1` valid | Body over 1 KiB | 413 | Pending | [ ] |
| Reverse Proxy HTTPS | Repeat all rows above | Same mutations | Same results | Pending | [ ] |

## Sensitive-data and lifecycle observations

| Check | Observation | Pass |
| --- | --- | --- |
| `Cache-Control: no-store` on every response | Pending | [ ] |
| No session material in URL, response, logs, policy, storage, or child arguments | Pending | [ ] |
| Package account remains the CGI process identity | Pending | [ ] |
| No Web Station, listener, root helper, nginx/global DSM modification, or new dependency | Pending | [ ] |
| Reverse Proxy rule and local policy removed during reset | Pending | [ ] |

## Decision

Pending. If any matrix row cannot meet the expected result, record the exact limitation and return the runtime decision to P1-002; do not replace the selected custom-CGI topology with an undocumented workaround.
