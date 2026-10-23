# Existing packages with File Station integration

Research date: 2026-10-23. No package installed, no downloaded code executed, and no NAS accessed. This evidence informs the conditional [Phase 3 File Station integration](phases/03-file-station-integration.md); it is not an implementation claim.

## Finding

**There is a concrete, undocumented package configuration mechanism named `fb_extern`.** I verified it in official legacy Text Editor and PDF Viewer SPKs downloaded from Synology. This is materially stronger evidence than the public developer guide alone provided.

These examples support investigating a package-owned context menu registration. They do not establish that arbitrary third-party registrations work on DSM 7.4, or that Synology supports this interface for independent developers. No Properties-tab equivalent was found.

The earlier recommendation to obtain Synology confirmation before any further work was unnecessarily restrictive. Confirmation is useful for determining official support; inspecting this mechanism and validating compatibility can proceed independently. The [roadmap](phases/README.md) keeps it optional and experimental, so it cannot block the standalone package.

## Packages examined

| Package/project | Integration found | Evidence and limits |
| --- | --- | --- |
| **Synology Text Editor** | “Open with Text Editor” in File Station; selected-file callback launches its application | Official help describes the action. Downloaded 1.2.3-0029 package contains `fb_extern` and corresponding JavaScript callbacks. Primary artifact verified; legacy version. [1][2] |
| **Synology PDF Viewer** | Package config declares a File Station menu action with selection and launch callbacks | Downloaded 1.2.3-0095 package contains `fb_extern`, including menu placement fields. Primary artifact verified; legacy version. [3] |
| **Synology Universal Viewer** | File previews from File Station | File Station 1.4.3 release notes describe integration. Current product precedent, but not evidence that `fb_extern` implements it or is available to third parties. [4] |
| **synOCR** | Opens File Station at a folder from its own UI, and uses File Station APIs for browsing | Source inspected. Its advertised “File Station integration” does not provide the requested File Station-to-package menu example. [5] |
| **SynoCommunity packages** | No relevant menu registration found in the inspected packaging source tree | Searched extracted `spk`, `mk`, and `docs` files. Most File Station references concern granting file permissions or ordinary browsing. This is a bounded negative result, not proof that no community package has ever implemented it. [6] |

## Text Editor: strongest example

In the official **TextEditor-x86_64-1.2.3-0029.spk**:

- `package.tgz` contains `ui/config` and `ui/TextEditor.js`.
- `ui/config` places `fb_extern` under the `SYNO.TextEditor.Application` application entry.
- Its registration names `SYNO.TextEditor.Utils.checkFn` and `SYNO.TextEditor.Utils.launchFn`.
- Other fields include `multiple`, `dir`, `text`, `icon`, `iconCls`, `dbclick`, and `vfs`.
- The launch function receives selected records, reads the first record’s `path`, and calls `SYNO.SDS.AppLaunch` with an `fb_recs` launch parameter.
- The check function examines application privilege, directory status, and the filename extension. These are client-side eligibility checks; a MediaInfo backend must independently authorize every request.

I also read all lifecycle scripts in this legacy package. They do not patch File Station’s scripts or bootstrap files. Its post-install script migrates Text Editor user preferences. Combined with the application-owned declaration, this supports the inference that File Station discovers the action through registered application metadata. The actual File Station registry reader has not yet been inspected, so the discovery logic and any allowlist remain unknown. [2]

A 2022 forum post independently shows the same registration in Text Editor **1.2.4-0245**. That is secondary corroboration of a later version, not current DSM 7.4 runtime proof. [7]

## PDF Viewer: a second registration

In **PDFViewer-x86_64-1.2.3-0095.spk**, `package.tgz` contains `ui/config`. The `SYNO.SDS.PDFViewer.Application` entry includes:

- `fb_extern` with `multiple: true` and `dir: false`;
- `checkFn` and `launchFn` pointing to its `FBExt` component;
- `inMenu` set to `gc_view`, `menuPriority` set to `-1`, and a translated action label.

This corroborates a shared registration structure rather than a Text Editor-only callback convention. The precise behavior of grouping, priority, and multi-selection must still be checked in the DSM 7.4 consumer. [3]

## What this changes for MediaInfo

The candidate approach is now:

1. Register a package-owned DSM application and its JavaScript bundle.
2. Investigate an `fb_extern` declaration that points to our own eligibility and launch callbacks.
3. Limit eligibility to supported regular media files, initially one selected file; avoid claiming the double-click action.
4. Pass the selection into the package’s own DSM window and request analysis through the authenticated backend described in the feasibility report.

This approach may satisfy the context menu requirement without replacing DSM/File Station files. **That remains a technical hypothesis until tested with a third-party package on DSM 7.4.** Do not copy or redistribute Synology’s application implementation; the proposed package would supply its own code.

An extra Properties tab remains a separate question. `fb_extern` evidence concerns file actions and application launch, not modification of the Properties dialog.

## Current-version limitation

I downloaded Text Editor **1.2.6-0260**, Text Editor **1.2.4-0245**, PDF Viewer **1.3.1-1233**, and File Station **1.4.4-2221** from Synology’s CDN. Unlike the older inspected artifacts, they were not readable by the ordinary tar parser and had Synology-format binary headers. I did not decrypt these archives. Synacktiv documents the encrypted SPK format. [8]

Consequently, this investigation does **not** claim to have inspected the current Text Editor callbacks or File Station’s DSM 7.4 registration consumer. File Station’s release notes confirm that 1.4.4 targets DSM 7.4. [4]

The most useful remaining read-only check is to inspect the installed packages’ `ui/config` files and the File Station code consuming `fb_extern` on a DSM 7.4 test system, or obtain equivalent extracted artifacts. Determine:

- whether `fb_extern` still exists and whether third-party applications are eligible;
- how callback bundles load and what argument/record shapes are passed;
- whether `path` is share-relative, physical, or virtual for each selection type;
- how enable/disable, install/uninstall, application privileges, translations, and File Station standalone mode affect the action;
- whether no-double-click registration is possible without changing existing media preview behavior.

## Reproducibility

Legacy packages were retrieved from the direct HTTPS URLs linked below, read as tar archives, and their nested `package.tgz` members inspected without installing them. No proprietary package code has been added to this workspace.

| Inspected artifact | SHA-256 |
| --- | --- |
| TextEditor-x86_64-1.2.3-0029.spk | `8d625ba9ab5ae374961382af5580ee65d2e078ba519eced5962ed5985af91aba` |
| PDFViewer-x86_64-1.2.3-0095.spk | `ae6882d296fa0b72df82998e2b16f99c9783f91060213fff73aae0f37b5494c5` |

Local scratch extraction: `/tmp/mediainfo-official-package-research`. synOCR source extraction: `/tmp/mediainfo-synocr-research`. SynoCommunity source extraction: `/tmp/mediainfo-package-research`. These are disposable research files.

Source revisions observed during this search:

- synOCR: `a5e8724c122e7133e0d26573cd6dd83cd9c96257`, dated 2026-10-17.
- SynoCommunity/spksrc: `0fb738cc65ddffbaf7ea82cc720d7187633052b5`, dated 2026-10-22.

The SynoCommunity search covered 1,480 extracted files under `spk`, `mk`, and `docs`, with a 1 MB per-file extraction limit. Search terms included `fb_extern`, File Station, FileStation, file handlers, file associations, and context menus. Bundled upstream applications, history, external source downloads, and other package repositories were not exhaustively searched.

## References

1. [Synology Text Editor help](https://kb.synology.com/th-th/DSM/help/TextEditor/text_editor?version=7): describes opening a selected file from File Station’s context menu.
2. [Official Text Editor 1.2.3-0029 SPK](https://global.synologydownload.com/download/Package/spk/TextEditor/1.2.3-0029/TextEditor-x86_64-1.2.3-0029.spk), reached through [Synology’s Text Editor archive](https://archive.synology.com/download/Package/TextEditor).
3. [Official PDF Viewer 1.2.3-0095 SPK](https://global.synologydownload.com/download/Package/spk/PDFViewer/1.2.3-0095/PDFViewer-x86_64-1.2.3-0095.spk) and [archive index](https://archive.synology.com/download/Package/PDFViewer/1.2.3-0095).
4. [File Station release notes](https://www.synology.com/en-us/releaseNote/FileStation?model=DS216%2B&os=DSM&version=7_x_series): Universal Viewer integration in 1.4.3 and DSM 7.4 requirement for 1.4.4.
5. [synOCR File Station helper source](https://github.com/geimist/synOCR/blob/a5e8724c122e7133e0d26573cd6dd83cd9c96257/APP/ui/template/synocr-filestation.js): app launch and open-path handling; [project README](https://github.com/geimist/synOCR).
6. [SynoCommunity/spksrc source revision](https://github.com/SynoCommunity/spksrc/tree/0fb738cc65ddffbaf7ea82cc720d7187633052b5).
7. [Text Editor 1.2.4 configuration posted in Synology-Forum.de](https://www.synology-forum.de/threads/text-editor-das-laden-der-voreinstellungen-ist-fehlgeschlagen.121211/post-1004179). Secondary discovery lead; core findings above were checked against original older artifacts.
8. [Synacktiv: Extraction of Synology encrypted archives](https://www.synacktiv.com/en/publications/extraction-of-synology-encrypted-archives-pwn2own-ireland-2024.html).
