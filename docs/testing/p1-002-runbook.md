# P1-002 — package request runtime proof

This runbook exercises the `mediainfo-poc` disposable prototype on `DSM-7.4-Test`. It selects a request runtime without adding an installed runtime dependency. The current prototype is a static `linux/amd64` Go FastCGI executable for the `v1000` VM, exposed through DSM's documented Web Config resource. It is a P1-002 probe only: it does not read files, invoke MediaInfo, accept request bodies, or establish the P1-003 authentication guarantees.

The initial PHP investigation is mandatory. A bundled Go fallback is selected only when the observations show that PHP does not provide a documented, package-owned request SAPI under the package account. An installed CLI binary, DSM's own PHP pool, or PHP's development server does not pass this gate.

## 1. Build the prototype

On the build host, run:

```bash
./scripts/test.sh
./scripts/build-spk.sh
file dist/stage/package/bin/runtime-probe
ldd dist/stage/package/bin/runtime-probe || true
sha256sum dist/mediainfo-poc-0.1.0-1-v1000.spk
```

The build uses the official `golang:1.25.4-bookworm` image, Go standard-library imports only, and `CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOAMD64=v1`. The NAS needs no Go compiler. Retain the source revision and SPK SHA-256 in the evidence report.

## 2. Refresh the target and assess PHP

Record DSM update, File Station version, package architecture, and installed-package baseline before installing the prototype. On an administrator guest shell, collect the following sanitized output. Do not include cookies, passwords, serial numbers, or package contents.

```sh
date -Iseconds
cat /proc/syno_platform
cat /proc/syno_cpu_arch
uname -m
ls -l /usr/bin/php* /usr/sbin/php* 2>/dev/null || true
for php in /usr/bin/php /usr/sbin/php /var/packages/*/target/bin/php; do
    [ -x "$php" ] || continue
    printf '\n== %s ==\n' "$php"
    "$php" -v
    "$php" -m
    "$php" --ini 2>/dev/null || true
done
find /usr -maxdepth 4 -type f -name '*php*fpm*' -perm -111 2>/dev/null
test -x /usr/syno/synoman/webman/modules/authenticate.cgi && echo helper-present
```

Also record the applicable published Synology documentation. As of 27 October 2026, Context7 supplied no authoritative Synology developer material for this route; use the primary [Web Config](https://help.synology.com/developer-guide/resource_acquisition/web_config.html), [Privilege Config](https://help.synology.com/developer-guide/privilege/privilege_config.html), and [Application Authentication](https://help.synology.com/developer-guide/integrate_dsm/web_authentication.html) pages.

Mark PHP **Pass** only if a package-owned CGI/FastCGI execution path, its package-account UID, usable configuration, and a documented support basis are all demonstrated. Otherwise record each failed gate and select Go.

## 3. Install and prove the selected route

In Package Center, manually install the generated SPK. Inspect the package details and installation result, then start it. Do not accept any new agreement without separately recording it.

For Go, capture these administrator-shell checks after startup:

```sh
synopkg status mediainfo-poc
id mediainfo-poc
ps -o user=,group=,pid=,args= -C runtime-probe
ls -ld /var/packages/mediainfo-poc /var/packages/mediainfo-poc/var
ls -l /var/packages/mediainfo-poc/var/runtime.sock
ss -lx | grep runtime.sock || true
ss -ltnp | grep runtime-probe || true
```

With an authenticated DSM browser session, request `http://NAS:5000/mediainfo-poc/runtime`. The expected success response is controlled JSON with `runtime`, `revision`, numeric `uid` and `gid`, and the DSM-authenticated `user`. Never copy a session cookie into an evidence report.

Repeat after stopping, starting, and restarting the package. After an uninstall/reinstall cycle, verify the route, socket, PID file, and process are absent after uninstall and return after reinstall. Record the browser-visible result and sanitized shell output.

## 4. P1-003 boundary

This probe only calls DSM's authentication helper with an allowlisted CGI environment and fails closed on helper failure or timeout. The package nginx fragment overwrites the backend-only DSM server-address parameter; P1-003 must separately prove valid/expired/logged-out sessions, application privileges, CSRF, origins, direct HTTPS/proxy behavior, request-size handling, and header forgery resistance. Do not interpret the P1-002 route as authorization to inspect any file.
