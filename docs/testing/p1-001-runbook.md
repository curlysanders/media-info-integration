# P1-001 — DSM test VM setup and reset

This procedure is for the disposable `DSM-7.4-Test` libvirt domain. It uses a DSM 7.4.1 guest presented as a DS1621+ by Arc Loader. Record the actual guest versions at each run; the emulated model and package architecture do not establish physical Synology compatibility. Enter dated observations in the [P1-001 evidence report](p1-001-evidence-2026-10-26.md), keeping expected results separate from observations. The detailed host handoff remains local to the operator.

## 1. Record the target before changing it

1. On the host, record the domain state, CPU and memory allocation, attached disk order, NAT lease, and the source paths with `virsh -c qemu:///system dominfo DSM-7.4-Test`, `domblklist DSM-7.4-Test --details`, and `domifaddr DSM-7.4-Test --source lease`. Do not commit the domain UUID or MAC address.
2. In DSM **Control Panel → Info Center → General**, record the complete DSM build/update, presented model, CPU string, core count, and memory. In **Info Center → Storage** and **Storage Manager**, record disk capacity, Volume 1 capacity/free space, filesystem, and storage-pool layout. Label any value seen only in the handoff as unverified.
3. In Package Center, record the installed File Station and SynoCommunity MediaInfo package versions, MediaInfo architecture build, and package source. If MediaInfo is absent, record that before installing the tested SynoCommunity build. Do not accept terms or install a package without the operator completing any required agreement.
4. When guest shell access is available, record the package architecture reported by DSM, then run `/var/packages/mediainfo/target/bin/mediainfo --Version` and record the executable version and exit status. A package listing alone does not verify the executable. Record any shell access method without storing credentials.
5. Mark prototype package version/revision and package-account identity **pending** until the prototype exists. Recheck the complete inventory after any DSM, File Station, MediaInfo, loader, or prototype update.

## 2. Check the provisional reset point

The existing 23 October host backup contains `domain.xml`, `arc.img`, and `DSM-7.4-Test.qcow2`; the running guest has changed since that copy. Before relying on it, inspect file presence/size, `qemu-img info` for both disks, the qcow2 backing chain, and `qemu-img check` on the backup image with sufficient host privileges. Record command exit codes and observations. These checks do not prove the backup is restorable. Keep the old backup intact; any new cold backup must use a dated destination, a fully shut-down VM, and both disks from the same point in time. For an isolated restore, use new disk paths, domain identity, and network settings; boot the copy before treating it as a reset point.

Until a separate restore boots and its data is checked, treat the old backup as a provisional recovery aid. Reset ordinary tests by removing only the named test resources below.

## 3. Create the ACL test area

Use a new share named `mediainfo-p1-test` on Volume 1 and ordinary local users `mediainfo-p1-a` and `mediainfo-p1-b`. Give neither user administrator membership. Record their effective group membership and any inherited permissions. Use only non-sensitive generated media. No package account receives a broad share grant.

Create `shared`, `a-only`, `b-only`, and `explicit-deny` folders. In File Station **Properties → Permission**, record each explicit or inherited Allow/Deny entry, the principal, the permissions, and its application scope. Synology documents this editor and its inheritance display in [File Station permissions](https://kb.synology.com/en-global/DSM/help/FileStation/privilege?version=7). Set the intended matrix below, then verify it in separate user sessions by opening or downloading each fixture. A folder being visible is not proof that its file is readable.

| Folder | User A | User B | Purpose |
| --- | --- | --- | --- |
| `shared` | Read | Read | Baseline allow |
| `a-only` | Read | Deny | Contrasting access |
| `b-only` | Deny | Read | Contrasting access |
| `explicit-deny` | Read | Explicit deny | Test Deny against an otherwise applicable Allow |

For the last row, create and record an applicable allow for user B through a group or parent, then add an explicit Deny for user B on the test folder/file. Verify the effective result; do not infer it from the editor alone. Record both session outcomes for every row, including unexpected access.

## 4. Prepare media and negative fixtures

Prepare a short MP4 and a short MKV from owned or freely usable source material. Each must contain identifiable video and audio tracks. Prefer differing codecs or track arrangements; subtitles, chapters, and descriptive tags are useful where available. Host-side `ffmpeg`/`ffprobe` may prepare and inspect fixtures; neither is an installed package runtime dependency. Record the exact commands or source and any tool versions used, then capture container, duration, video/audio codec names, additional tracks, tags, byte size, and SHA-256 checksum. Keep the media out of Git.

The following host-side recipe creates the initial synthetic fixtures. It was exercised with FFmpeg/FFprobe 6.1.1-3ubuntu5. The `$(not-executed)` string is literal tag data inside single quotes.

```bash
mkdir -p /tmp/mediainfo-p1-001-fixtures
ffmpeg -hide_banner -loglevel error -f lavfi -i testsrc2=size=320x180:rate=24:duration=3 -f lavfi -i sine=frequency=440:sample_rate=48000:duration=3 -c:v libx264 -preset ultrafast -pix_fmt yuv420p -c:a aac -b:a 96k -metadata title='P1-001 MP4 codec fixture' -shortest -y /tmp/mediainfo-p1-001-fixtures/normal-h264-aac.mp4
ffmpeg -hide_banner -loglevel error -f lavfi -i testsrc2=size=320x180:rate=24:duration=3 -f lavfi -i sine=frequency=660:sample_rate=48000:duration=3 -c:v libx265 -preset ultrafast -pix_fmt yuv420p -c:a libopus -b:a 96k -metadata title='P1-001 MKV codec fixture' -shortest -y /tmp/mediainfo-p1-001-fixtures/normal-hevc-opus.mkv
printf '1\n00:00:00,000 --> 00:00:02,500\nP1-001 subtitle track\n' > /tmp/mediainfo-p1-001-fixtures/subtitles.srt
ffmpeg -hide_banner -loglevel error -i /tmp/mediainfo-p1-001-fixtures/normal-hevc-opus.mkv -i /tmp/mediainfo-p1-001-fixtures/subtitles.srt -map 0 -map 1 -c copy -c:s srt -metadata:s:s:0 language=eng -y /tmp/mediainfo-p1-001-fixtures/normal-hevc-opus-srt.mkv
ffmpeg -hide_banner -loglevel error -i /tmp/mediainfo-p1-001-fixtures/normal-h264-aac.mp4 -map 0 -c copy -metadata title='<svg onload="alert(1)">& $(not-executed) 🧪' -metadata comment=$'line one\nline two; --Output=HTML' -y /tmp/mediainfo-p1-001-fixtures/hostile-tags.mp4
printf 'P1-001 intentionally malformed media\000\377\001' > /tmp/mediainfo-p1-001-fixtures/malformed.bin
head -c 512 /tmp/mediainfo-p1-001-fixtures/normal-h264-aac.mp4 > /tmp/mediainfo-p1-001-fixtures/truncated.mp4
```

The final MKV is `normal-hevc-opus-srt.mkv`; `normal-hevc-opus.mkv` is an intermediate. Copy the five final media/negative files into the guest share when it exists. Make a known-readable clip copy in each denied ACL location and use filenames with spaces and option-like prefixes for later invocation tests. Recalculate byte sizes and checksums after transfer. Do not use private media or embed credentials in metadata.

Place the normal clips in `shared` and ACL copies in the appropriate folders. Use separate sessions to verify the matrix. P1-001 requires representative content, not a large-file benchmark. Later tickets must still define and exercise numeric size and resource limits; report large-media latency and storage behavior as unmeasured.

## 5. Reset and preserve evidence

Before cleanup, record expected versus observed results, timestamps, exact versions, fixture hashes and sizes, ACL entries, and any limitations in the evidence report. Do not capture passwords, cookies, session tokens, serial numbers, media payloads, or proprietary package contents.

Remove only test-created prototype registrations/packages, media, share, and test users/groups, in that order. Verify removal in Package Center, File Station, and Control Panel. Remove MediaInfo only if this procedure installed it solely for this test. Preserve unrelated packages and the existing host backup. Record each removal result and any stale registration or leftover files. If a guest-wide reset is required, use an isolated, validated restore rather than overwriting disks attached to the running domain.
