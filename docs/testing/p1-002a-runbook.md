# P1-002A — documented identity hand-off research gate

This runbook determines whether the disposable `DSM-7.4-Test` VM has a **documented** way for a standalone third-party package to derive the DSM session identity. It does not authorize file access, MediaInfo execution, or P1-003.

## Source gate

Record the research date, DSM build, package revision, and the following primary sources before testing:

- [Desktop Application](https://help.synology.com/developer-guide/integrate_dsm/desktopapp.html) and [INFO optional fields](https://help.synology.com/developer-guide/synology_package/INFO_optional_fields.html), which document `dsmuidir` and its managed `/usr/syno/synoman/webman/3rdparty/<identifier>/` link.
- [Application Authentication](https://help.synology.com/developer-guide/integrate_dsm/web_authentication.html), which recommends invoking `authenticate.cgi` directly from a package custom CGI because its request environment is set automatically.
- [Privilege Config](https://help.synology.com/developer-guide/privilege/privilege_config.html), which defines the package-account model.
- [Web Config](https://help.synology.com/developer-guide/resource_acquisition/web_config.html) and [Web Service](https://help.synology.com/developer-guide/resource_acquisition/web_service.html), to record why the prior FastCGI and Web Station candidates are not the selected topology.

Context7 was consulted first on 27 October 2026. Its available Synology DSM entry contained no matching developer-guide material, so this gate relies on the dated primary sources above. The published developer guide identifies itself as DSM 7.2.2; it is evidence for a DSM 7.4 test candidate, not a DSM 7.4 compatibility guarantee.

## Candidate and package probe

`mediainfo-poc` `0.1.0-0011` replaces the Web Config/FastCGI listener with `ui/runtime.cgi`. `dsmuidir="ui"` asks DSM to create the package-owned third-party UI link. The expected probe URL is:

```text
https://<nas>:5001/webman/3rdparty/mediainfo-poc/runtime.cgi
```

The CGI invokes `/usr/syno/synoman/webman/modules/authenticate.cgi` directly, with only bounded CGI values (`HTTP_COOKIE`, remote/server address, request method and URI, protocol, host, and related CGI fields). It never accepts a username or forwarded identity as proof, uses no shell, emits `Cache-Control: no-store`, and exposes no TCP listener or persistent process.

Build and record the exact source revision and digest:

```bash
./scripts/test.sh
./scripts/build-spk.sh
file dist/stage/package/ui/runtime.cgi
ldd dist/stage/package/ui/runtime.cgi || true
sha256sum dist/mediainfo-poc-0.1.0-0011-v1000.spk
```

Install the SPK manually in Package Center. Do not add Web Station or another runtime, alter DSM/nginx files, accept a root privilege declaration, or add a helper/broker. Inspect the managed `dsmuidir` link and the CGI process identity with an administrator guest shell; redact all cookies and session values.

```sh
synopkg status mediainfo-poc
id mediainfo-poc
readlink /usr/syno/synoman/webman/3rdparty/mediainfo-poc
ls -l /var/packages/mediainfo-poc/target/ui/runtime.cgi
```

## Sanitized request matrix

Use separate non-administrator DSM browser sessions. Record timestamp, request class, HTTP status, response headers, response body shape, execution UID/GID, and observed identity; never record cookies, passwords, serial numbers, or raw request headers.

| Case | Expected result |
| --- | --- |
| Valid active session | `200`, controlled JSON, identity produced by DSM helper |
| No session | `401`, controlled JSON, `Cache-Control: no-store` |
| Logged-out session | `401` after DSM logout |
| Expired session | `401` after normal expiry or a documented expiry test |
| Forged username / forwarded header / query value | Cannot select or alter identity; only a valid DSM session may return `200` |
| Malformed or copied fake cookie | `401`, no session value echoed or logged |
| POST | `405`, controlled JSON |

After every response, verify the request path is a DSM-managed package link rather than a manually-added nginx route, no Web Station/PHP package was installed, no listener/process remains after the request, and the observed CGI account satisfies the package-account security invariant. Test package stop/start, upgrade, uninstall, and reinstall; the managed UI link and executable must disappear on uninstall.

## Decision rule

Pass P1-002A only when the source gate documents the full package-to-custom-CGI route and the live matrix proves the authenticated DSM identity reaches the package without an excluded dependency or privilege. Update P1-002 with the selected runtime, route, revision, account, and unresolved P1-003 assumptions.

If the managed link does not execute the CGI, the process is not acceptably package-scoped, the authenticated session is rejected, or any excluded component is needed, publish **unavailable** for DSM `7.4.1-90080 Update 1`. Keep P1-002 and P1-003–P1-006 blocked; do not restore the FastCGI route or implement a workaround.
