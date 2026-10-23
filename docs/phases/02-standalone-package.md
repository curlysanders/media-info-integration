# Phase 2 — Standalone DSM package

## Objective

Ship an open-source DSM 7.4 SPK that lets an authorized user choose one media file and view MediaInfo in a package-owned native-looking DSM window. This is the first usable product and has no dependency on a File Station menu.

## Prerequisites

- Phase 1 has selected and validated the package runtime and authorization approach.
- SynoCommunity MediaInfo compatibility has been verified for each supported NAS architecture.
- The documented DSM application-window route has opened successfully on the target DSM build.

## Work

- Build a package-owned DSM application window using the documented app integration surface. Use plain TypeScript for application logic and only the minimal DSM-native UI adapter needed by that route; do not add a separately installed frontend runtime or a large framework.
- Provide a small file picker backed by documented File Station APIs. It works independently of a context-menu selection and displays only files the user is authorized to select.
- Keep the request handler plain PHP without Composer or Symfony when Phase 1 validates built-in PHP; otherwise ship the selected Go executable. The runtime and frontend must not add an installed NAS dependency.
- Render General, Video, Audio, Text, and Chapters sections, with a copyable raw report. Use text rendering for all values, accessible loading and error states, and one selected file at a time.
- Run under the package account. Re-check authentication, CSRF, authorization, request limits, and MediaInfo availability on every analysis request. Start without a database, shared cache, batch processing, or a persistent media index.
- Declare and verify the `mediainfo` package dependency. Report a clear, non-sensitive error when it is missing, disabled, incompatible, or updated.
- Build with a pinned Synology Package Toolkit and reproducible CI inputs. Publish manual-install SPKs and checksums through GitHub Releases for only tested DSM versions and CPU architectures. Describe the evidence for reproducibility; do not claim byte-identical builds without proving it.

## Deliverables

- Source, package metadata, build instructions, and a dependency/version compatibility matrix.
- An installable SPK for each tested architecture, plus checksum and release notes.
- A standalone DSM window with selection, analysis, display, error handling, and cancellation.
- Install, upgrade, restart, dependency removal/update, and uninstall evidence.

## Acceptance checklist

- [ ] An authorized user can install the package, select a permitted media file, and see its MediaInfo result.
- [ ] A denied selection, malformed file, missing dependency, timeout, and cancellation produce safe and understandable errors.
- [ ] The package has no root process, system-file modification, broad share grant, secret logging, or retained media copies.
- [ ] Upgrade and uninstall remove package-owned state and leave DSM/File Station unchanged.
- [ ] Build and installation instructions name the exact supported DSM and architecture matrix.

## Blockers and exclusions

If Phase 1 cannot establish secure authorization to the original file, do not release a package that analyzes paths accessible only to a privileged service. Do not add File Station menu integration, Properties tabs, batch jobs, a database, or a persistent cache in this phase.
