# Changelog

All notable changes to the Tomba CLI will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.1.3] - 2026-04-02

### Added

- **AI Chat command** (`tomba chat`) — natural language queries powered by OpenAI with Tomba function calling. Chains tools like `domain_search -> email_finder -> email_verifier` automatically.
- **Chat bulk mode** (`tomba chat --file`) — process CSV files through AI chat with progress bar and markdown summary table.
- **Bulk CSV command** (`tomba bulk`) — batch process CSV files with auto column mapping. Supports `enrich`, `verify`, `finder`, and `search` operations.
- **ClawHub Skills** (`tomba skill`) — reusable multi-step Tomba API workflows. 8 built-in skills: `find-email`, `company-emails`, `enrich-verify`, `company-intel`, `lead-gen`, `linkedin-intel`, `domain-check`, `author-verify`.
- Skill management: `skill list`, `skill run`, `skill info`, `skill create`, `skill import`, `skill export`, `skill remove`.
- Custom skill YAML format with `{{input.name}}` and `{{steps.id.field}}` template engine for chaining step results.
- **Text output mode** — human-readable formatted output is now the default (`--json` to opt in to JSON).
- Unified `output.Render()` function replaces repetitive output code across all commands.
- Rich text formatters for all commands with full field coverage: contact details, social profiles, verification status, sources, department breakdowns, progress bars, and more.
- **Install scripts** for Linux/macOS (`install.sh`) and Windows (`install.ps1`) with Tomba ASCII art banner.
- SVG demo recordings for all 32 commands using termtosvg.
- Shell completions: bash, zsh, fish.
- Man pages.

### Changed

- Default output is now text instead of JSON. Use `--json` or `-j` to get JSON output.
- All command output refactored to use `output.Render()` — removed ~240 lines of duplicated output code.
- `formatReveal` updated to support new API response format (`data.companies[]`, pagination, active filters, full company details).
- `formatSimilar` updated to support `website_url`, `industries` fields and `meta` pagination.
- `formatTechnology` updated to support flat `data[]` array and `categories` as array of objects.
- `formatPhoneFinder` updated to support single phone object, array of phones, and nested `phones[]` formats.
- `formatStatus` updated to handle both flat and `data`-wrapped responses.

### Fixed

- Status command not displaying output when API returns flat JSON without `data` wrapper.
- Phone finder not rendering when API returns single object instead of `phones[]` array.
- Reveal command missing companies when API uses `data.companies[]` structure.
- Similar command using wrong field names (`host` instead of `website_url`).
- Technology command expecting `data.technologies[]` when API returns flat `data[]`.

## [1.1.2] - 2026-03-10

### Changed

- Update GitHub Actions to use `GH_PAT` for `GITHUB_TOKEN`.
- Add directory configuration to Scoop bucket.

## [1.1.1] - 2026-03-10

### Fixed

- Remove icon from Snapcraft configuration.

## [1.1.0] - 2026-03-10

### Changed

- Enhance Goreleaser workflow configuration.
- Remove separate Snapcraft workflow (consolidated into Goreleaser).

## [1.0.9] - 2026-03-10

### Changed

- Enhance Snapcraft configuration.
- Remove obsolete Go workflow configuration.

### Added

- Installation guide for Snap from dist folder.

## [1.0.8] - 2026-03-10

### Added

- Snapcraft build and publish workflow.
- `phone-validator` command — validate phone numbers.
- `phone-finder` command — search phone numbers by email, domain, or LinkedIn URL.
- `reveal` command — search companies using natural language or structured filters.
- `similar` command — retrieve domains similar to a specific domain.
- `technology` command — discover technologies detected for a domain.
- `whoami` command — print current account information.
- `enrich-mobile` flag for finder, enrich, and linkedin commands.
- New HTTP server routes for all new commands.

### Changed

- Update Tomba SDK dependency to v1.0.7.
- Improve email handling with FinderData extraction.
- Update GoReleaser action to version 6.

### Fixed

- Email verification utility description.
- Time description in CLI messages (minutes to seconds).

## [1.0.7] - 2024-08-07

### Added

- `sources` command — find where an email address has been found on the web.
- Tomba systemd service configuration file.
- Post-install script for enabling systemd service.
- Shell completion generation (bash, zsh, fish).
- Man page generation.

### Changed

- Update Goreleaser configuration with signing and package contents.
- Update Tomba SDK dependency to v1.0.3.

## [1.0.6] - 2024-03-12

### Added

- Docker support.
- ChatGPT plugin integration.

### Fixed

- OpenAPI call handling.
- Total emails count on domain search.
- Slack response formatting.

## [1.0.5] - 2023-08-08

### Added

- `finder` command — email finder by domain and name.
- Slack integration with slash commands (`/search`, `/enrich`, `/author`, `/linkedin`, `/checker`).
- Slack manifest and response builders.
- Homebrew and Scoop installation support.
- Domain search flags (page, limit, department).

## [1.0.4] - 2023-07-28

### Added

- New HTTP route endpoints.
- New CLI commands.

### Fixed

- Login flow improvements.

## [1.0.3] - 2023-07-27

### Added

- `http` command — runs a reverse proxy HTTP server.
- HTTP endpoints for all Tomba API operations.
- `--port` flag for configuring HTTP server port.
- Homebrew and Scoop package manager support.

## [1.0.2] - 2023-07-26

### Added

- Snap Store badge in README.

### Fixed

- Application name correction.
- Repository name in configuration.

## [1.0.1] - 2023-07-26

### Added

- Search parameters (page, limit).
- Multiple architecture builds (386, amd64, arm, arm64, ppc64).

### Fixed

- LinkedIn finder functionality.

### Changed

- Update dependencies.

## [1.0.0] - 2023-07-26

### Added

- Initial release of Tomba CLI.
- `login` / `logout` — authentication with API key and secret.
- `search` — domain search for email addresses.
- `enrich` — email enrichment with contact data.
- `verify` — email deliverability verification.
- `count` — email count per domain.
- `status` — domain webmail/disposable detection.
- `author` — article author email finder.
- `linkedin` — LinkedIn profile email finder.
- `usage` — monthly API request usage.
- `logs` — last 1,000 API request logs.
- JSON and YAML output formats with syntax highlighting.
- File output with `-o` flag.
- Snap package distribution.
- Goreleaser multi-platform builds (Linux, macOS, Windows).

[1.1.3]: https://github.com/tomba-io/tomba/compare/v1.1.2...HEAD
[1.1.2]: https://github.com/tomba-io/tomba/compare/v1.1.1...v1.1.2
[1.1.1]: https://github.com/tomba-io/tomba/compare/v1.1.0...v1.1.1
[1.1.0]: https://github.com/tomba-io/tomba/compare/v1.0.9...v1.1.0
[1.0.9]: https://github.com/tomba-io/tomba/compare/v1.0.8...v1.0.9
[1.0.8]: https://github.com/tomba-io/tomba/compare/v1.0.7...v1.0.8
[1.0.7]: https://github.com/tomba-io/tomba/compare/v1.0.6...v1.0.7
[1.0.6]: https://github.com/tomba-io/tomba/compare/v1.0.5...v1.0.6
[1.0.5]: https://github.com/tomba-io/tomba/compare/v1.0.4...v1.0.5
[1.0.4]: https://github.com/tomba-io/tomba/compare/v1.0.3...v1.0.4
[1.0.3]: https://github.com/tomba-io/tomba/compare/v1.0.2...v1.0.3
[1.0.2]: https://github.com/tomba-io/tomba/compare/v1.0.1...v1.0.2
[1.0.1]: https://github.com/tomba-io/tomba/compare/v1.0.0...v1.0.1
[1.0.0]: https://github.com/tomba-io/tomba/releases/tag/v1.0.0
