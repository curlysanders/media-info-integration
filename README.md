# MediaInfo for Synology DSM

This proof-of-concept is concluded; no DSM package will be developed further. The desired File Station **View MediaInfo** action is unavailable on the tested DSM 7.4/File Station target, while a standalone viewer is explicitly not the intended product experience. See the [final decision](docs/final-decision-2026-11-03.md).

| Document | Purpose |
| --- | --- |
| [Delivery roadmap](docs/phases/README.md) | Three phases, scope, prerequisites, acceptance checks, and blockers |
| [Feasibility investigation](docs/feasibility.md) | Supported Synology interfaces, security model, architecture options, and references |
| [Existing-package investigation](docs/package-integration-research.md) | Evidence for the undocumented legacy `fb_extern` mechanism and its limits |
| [Git workflow](docs/workflow.md) | Short-lived branches, signed commits, pull requests, and squash merging |
| [Contributor guidance](AGENTS.md) | Runtime, dependency, security, documentation, and DSM boundary rules |

The only planned externally installed runtime dependency is SynoCommunity MediaInfo. Built-in DSM PHP is assessed during Phase 1; the P1-002A disposable prototype contains a small bundled Go custom-CGI candidate because the prior FastCGI route could not preserve DSM session identity. No File Station Properties-tab hook has been found.
