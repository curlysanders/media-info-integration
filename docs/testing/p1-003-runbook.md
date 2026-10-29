# P1-003 — protected package routing proof

This procedure is for the disposable `DSM-7.4-Test` VM only. It validates the custom-CGI route selected by P1-002A; it does not select a file, call MediaInfo, or authorize filesystem access. Record actual target versions and sanitized observations in the dated evidence report before marking P1-003 complete.

## 1. Build an explicitly scoped test policy

The packaged `routing-policy.json` is a non-secret P1 test policy. It has one allowed user and an exact origin allowlist. It is not a Phase 2 authorization model. Create an untracked local copy with `test1` and the two exact HTTPS origins used in this run: the direct DSM HTTPS origin and the DSM Reverse Proxy HTTPS origin. Do not use wildcards, HTTP origins, `null`, request-derived hosts, passwords, session cookies, or DSM tokens.

```json
{
  "allowedUsers": ["test1"],
  "allowedOrigins": [
    "https://nas.example.test:5001",
    "https://mediainfo.example.test"
  ]
}
```

Build with the local policy path. The policy is copied inside the SPK next to `runtime.cgi`; it is never written by the CGI.

```bash
P1_ROUTING_POLICY=/absolute/path/p1-003-routing-policy.json ./scripts/test.sh
P1_ROUTING_POLICY=/absolute/path/p1-003-routing-policy.json ./scripts/build-spk.sh
```

Record the source revision, policy allowlist values, SPK digest, package revision, package UID/GID, DSM build, and both URLs. Do not record session material.

## 2. Configure the routes

Install or upgrade the SPK through Package Center on the disposable VM and verify that the DSM-managed package link targets the installed package UI directory. Use the documented DSM Reverse Proxy UI at **Control Panel → Login Portal → Advanced → Reverse Proxy** to create one disposable HTTPS source rule. Its HTTPS destination must be the same VM's DSM HTTPS service and its path must reach exactly:

```text
/webman/3rdparty/mediainfo-poc/runtime.cgi
```

Use a source hostname covered by a valid test certificate. Do not add custom identity headers, a Web Station dependency, a package nginx resource, a global nginx change, a listener, a root broker, or a trusted-proxy configuration for this ticket. The proxy is transport coverage only; its headers are not trusted by the CGI.

Confirm the direct URL and proxied URL both reach the managed package route. Record the source/destination protocol, hostname, port, path behavior, certificate result, and whether DSM preserves the expected POST request context. Delete the disposable rule during cleanup.

## 3. Run the request matrix

Use independent non-administrator DSM browser sessions for `test1` and `test2`. Obtain the DSM-generated `X-SYNO-TOKEN` only for the live request; never place it in shell history, a URL, evidence, source code, or persistent files. Send the exact probe document as JSON:

```json
{"operation":"route-probe"}
```

For each direct and proxied route, record only the case label, timestamp, HTTP status, response-header names/values, body shape, returned identity where successful, and package UID/GID. Redact all cookie and token values.

| Case | Expected result |
| --- | --- |
| `test1`, valid session/token, configured origin, valid JSON | `200`; controlled JSON for `test1`; `Cache-Control: no-store` |
| `test2`, valid session/token, configured origin | `403`; generic forbidden response |
| no session, forged/copy cookie, missing/invalid token, logged-out or expired session | `401`; generic authentication response |
| forged username, query value, `X-Forwarded-User`, `X-Forwarded-Host`, or `X-Forwarded-Proto` | Cannot alter identity or authorization result |
| missing/disallowed origin | `403` after valid DSM authentication |
| GET/PUT/other method | `405` |
| missing/wrong content type; empty, malformed, unknown-field, trailing JSON; malformed length | `400` |
| declared or actual body larger than 1 KiB | `413` |

For every response that reaches the browser as a package-controlled success or error, verify JSON with `Cache-Control: no-store`; it must not echo a token, cookie, supplied username, forwarded header, or request body. DSM `7.4.1-90080 Update 1` was observed to replace the package's `403` response on this managed route with a generic DSM error representation and to remove `Cache-Control`; record that behavior separately. The same DSM version can suppress the body of a package `401` after logout while retaining `Cache-Control: no-store`. These are acceptable only when the DSM response contains no identity, token, request, file, or other caller-specific information. Inspect the package process and any available sanitized logs to confirm no token/cookie is retained or passed to a child other than the bounded `authenticate.cgi` invocation. The browser-facing operation remains POST-only. Before calling the helper, the CGI requires `CONTENT_LENGTH`, reads exactly that declared size (maximum 1 KiB) from the live CGI input stream, rejects incomplete/malformed input, and supplies the helper the original bounded CGI method, content headers, and an in-memory copy of that body; neither body nor session material is persisted. Re-run the complete matrix to prove that the helper rejects missing, invalid, logged-out, and expired sessions.

## 4. Decision and cleanup

P1-003 passes only if the full matrix succeeds through both HTTPS routes and no sensitive material is retained. If proxying cannot preserve the DSM-authenticated custom-CGI flow, or the DSM token cannot reject cross-origin/malformed requests, update the P1-002 runtime decision instead of adding a workaround.

Remove the Reverse Proxy rule, local test policy, test package, and temporary users according to the P1-001 reset procedure. Preserve only sanitized evidence. The upstream Context7 DSM entry returned no custom-CGI authentication documentation on 27 October 2026; retain that limitation and cite the dated Synology Application Authentication guide used by P1-002A.
