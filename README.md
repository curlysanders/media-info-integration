# MediaInfo for Synology DSM

An open-source DSM 7.4 package is planned, not implemented. Phase 2 delivers a secure standalone MediaInfo viewer with a package-owned DSM window. Phase 3 may add an experimental File Station **View MediaInfo** action if DSM 7.4 accepts the observed legacy compatibility mechanism.

| Document | Purpose |
| --- | --- |
| [Delivery roadmap](docs/phases/README.md) | Three phases, scope, prerequisites, acceptance checks, and blockers |
| [Feasibility investigation](docs/feasibility.md) | Supported Synology interfaces, security model, architecture options, and references |
| [Existing-package investigation](docs/package-integration-research.md) | Evidence for the undocumented legacy `fb_extern` mechanism and its limits |
| [Git workflow](docs/workflow.md) | Short-lived branches, signed commits, pull requests, and squash merging |
| [Contributor guidance](AGENTS.md) | Runtime, dependency, security, documentation, and DSM boundary rules |

The only planned externally installed runtime dependency is SynoCommunity MediaInfo. Built-in DSM PHP will be assessed during Phase 1; a small bundled Go executable is the fallback if PHP cannot safely host the package route. No File Station Properties-tab hook has been found.
