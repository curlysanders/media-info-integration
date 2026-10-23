# Phase 3 — Optional File Station integration

## Objective

Offer a package-owned **View MediaInfo** action for one selected supported media file in File Station, opening the Phase 2 window with that selection. This phase proceeds only if Phase 1 proves the DSM 7.4 compatibility route.

## Prerequisites

- The standalone Phase 2 package is complete and remains usable without this feature.
- Phase 1 has documented a successful third-party `fb_extern` registration for the exact DSM and File Station versions to be supported.
- The handoff preserves the same backend authorization as direct use of the standalone application.

## Work

- Register the minimal package-owned `fb_extern` configuration and callbacks verified in Phase 1.
- Limit the action to one supported regular media file. Keep existing double-click/preview behavior untouched.
- Handoff a validated selection to the package window, then repeat server-side authentication and authorization before analysis. A visible menu item is never access control.
- Make the feature explicitly experimental and opt-in while the mechanism remains undocumented. Scope compatibility claims to exact verified DSM and File Station builds.
- Handle callback failures, disabled package state, dependency absence, upgrades, and uninstall without breaking File Station. The standalone app must continue to work when the integration is unavailable.

## Deliverables

- Version-scoped File Station compatibility notes and an opt-in setting with clear feature status.
- A one-file launch path that opens the Phase 2 window and produces the same authorized result as standalone selection.
- Lifecycle test evidence for install, enable/disable, upgrade, dependency change, and uninstall.

## Acceptance checklist

- [ ] The action appears only for an eligible single file on the validated DSM/File Station build.
- [ ] The selected path is handed off correctly and the backend independently rejects unauthorized or altered selections.
- [ ] Existing File Station preview and double-click behavior remains intact.
- [ ] Failure leaves File Station functional and users can use the standalone app.
- [ ] Disabling or uninstalling the package removes its registration without DSM or File Station changes.

## Blockers and exclusions

If DSM 7.4 rejects third-party registration or its behavior cannot be made safe and reversible, omit this feature from the standard package. Do not patch File Station files, inject into browser pages, decrypt or copy Synology implementation code, or add a Properties-window tab.
