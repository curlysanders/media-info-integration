# Delivery roadmap

This roadmap turns the feasibility work into a package in three stages. The target is DSM 7.4. Each stage records exact tested versions because Synology's published developer material does not establish every DSM 7.4 behavior.

| Phase | Outcome | Status |
| --- | --- | --- |
| [1. Proof of concept](01-proof-of-concept.md) | Evidence for a safe request runtime, per-user file authorization, MediaInfo analysis, and the `fb_extern` compatibility route | In progress (P1-001) |
| [2. Standalone package](02-standalone-package.md) | A usable open-source SPK with its own DSM window and file selection | Planned |
| [3. File Station integration](03-file-station-integration.md) | Optional experimental “View MediaInfo” menu action, if Phase 1 validates it | Conditional |

Phase 2 does not depend on Phase 3. A Properties-window tab is outside this roadmap because no supported or observed extension point has been found.

The [Phase 1 ticket backlog](01-tickets.md) breaks the proof of concept into ordered investigations with acceptance criteria and required evidence.

Supporting evidence is in the [feasibility report](../feasibility.md) and [existing-package investigation](../package-integration-research.md).
