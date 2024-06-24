# Changelog

The changelog consists of eight categories:
- **Added** - new functionality that brings value to users
- **Changed** - changes in existing functionality, performance and other types of improvements
- **Fixed** - bugfixes
- **Deprecated** - soon-to-be removed user-facing features
- **Removed** - earlier deprecated, now removed user-facing features 
- **Security** - fixing CVEs in the gateway or dependencies
- **Miscellaneous** - all other updates like build, release, CLI, etc.

See [keepachangelog.com](https://keepachangelog.com/en/1.1.0/) for more information.

## [Unreleased]

### Added

TBU

### Changed

TBU

### Fixed

TBU

### Deprecated

TBU

### Removed

TBU

### Security

TBU

### Miscellaneous

TBU

## [0.1.0-rc.1] (Jun 24, 2024)

The first major update with breaking changes to the language chat schemas 
and begging of work on instrumenting the gateway with OpenTelemetry.

### Added

- 🔧 Use github.com/EinStack/glide as module name to support go install cmd (@gernest)
- ✨🔧 Setup Open Telemetry Metrics and Traces (#237) (@gernest)
-  🔧 #221 Add B3 trace propagator (#242) (@gernest)
- 🔧 #241 Support overriding OTEL resource attributes (#243) (@gernest)
- 🔧 #248 Disable span and metrics by default (#254) (@gernest)
- 🔧 #220 Instrument API server with observability signals (#255) (@gernest)
- 🔧 #164 Make client connection pool configurable across all providers (#251) (@daesu)
- 🔧 Instrument gateway process (#256) (@gernest)
- 🔧 #262: adding connection pool for chat request and response (#271) (@tom-fitz)

### Changed

- 🔧 #238 Implements human-readable durations in config (#253) (@ppmdo)
- 🔧 #266: removing omitempty from response definition (#267) (@tom-fitz)

#### Breaking Changes

- 🔧 💥 #235: Extended the non-streaming chat error schema with new fields to give clients more context around the error (#236) (@roma-glushko)
- 💥 Convert all camelCase config fields to the snake_case in the provider configs (#260) (@roma-glushko)
- ✨💥 #153: Allow to pass multiple model-specific param overrides (#264) (@roma-glushko)

### Fixed

- 🐛 #217: Set build info correctly in Glide images (#218) (@roma-glushko)

### Security

- 🔒 Updated golang to 1.22.4 to address CVE-2024-24790 (#276) (@STAR-173)

### Miscellaneous

- 📝 Defined a way to manage EinStack Glide project (#234) (@roma-glushko)
- 👷 #219: Setup local telemetry stack with Jaeger, Grafana, VictoriaMetrics and OTEL Collector (#225) (@roma-glushko)
- 👷‍♂️ Added a new GH action to watch for glide activity stream (#239, #244) (@roma-glushko)
- ✨ Switched to the new docs (@roma-glushko)
- 🔧 #240: Automatically install air (#277, #270) (@ppmdo, @roma-glushko)

## [0.0.3-rc.2], [0.0.3] (Apr 17, 2024)

Final major improvements to streaming chat workflow. Fixed issues with Cohere streaming chat. 
Expanded and revisited Cohere params in config.

### Added

- 🔧 #195 #196: Set router ctx in stream chunks & handle end of stream in case of some errors (@roma-glushko)
- 🐛🔧 #197: Handle max_tokens & content_filtered finish reasons across OpenAI, Azure and Cohere (@roma-glushko)

### Changed

- 🔧 💥 #198: Expose more Cohere params & fixing validation of provider params in config (breaking change) (@roma-glushko)
- 🔧 #186: Rendering Durations in a human-friendly way (@roma-glushko)

### Fixed

- 🐛 #209: Embed Swagger specs into binary to fix panics caused by missing swagger.yaml file (@roma-glushko)
- 🐛 #200: Implemented a custom json per line stream reader to read Cohere chat streams correctly (@roma-glushko)

## [0.0.3-rc.1] (Apr 7th, 2024)

Bringing support for streaming chat in Glide.

### Added

- ✨Streaming Chat Workflow #149 #163 #161 (@roma-glushko)
- ✨Streaming Support for Azure OpenAI #173 (@mkrueger12)
- ✨Cohere Streaming Chat Support #171 (@mkrueger12)
- ✨Start counting token usage in Anthropic Chat #183 (@roma-glushko)
- ✨Handle unauthorized error in health tracker #170 (@roma-glushko)

### Fixed

- 🐛 Fix Anthropic API key header #183 (@roma-glushko)

### Security

-  🔓 Update crypto lib, golang, fiber #148 (@roma-glushko)

### Miscellaneous

-  🐛 Update README.md to fix helm chart location #167 (@arjunnair22)
- 🔧 Updated .go-version (@roma-glushko)
-  ✅ Covered the telemetry by tests #146 (@roma-glushko)
- 📝 Separate and list all supported capabilities per provider #190 (@roma-glushko)

## [0.0.2-rc.2], [0.0.2] (Feb 22nd, 2024)

### Added

- ✨ [Lang Chat Router] Ollama Support #142 (@mkrueger12)
- ✨ [Lang Chat Router] AWS Bedrock Support #131 (@mkrueger12)

### Miscellaneous

- 👷 Fixing the dockerhub authorization step in the release workflow #155 (@roma-glushko)
- ♻️ Moved specific provider schemas closer to provider's packages #151 (@roma-glushko)

## [0.0.2-rc.1] (Feb 12th, 2024)

### Added

- ✨ Allow to load dotenv files #117 (@roma-glushko)

### Changed

- ✨👷 Support for Windows #91 (@roma-glushko)
- 👷 Build Glide for OpenBSD and ppc65le, s390x, riscv64 architectures #139 (@roma-glushko)

### Miscellaneous

- 👷 Release binaries to Snapcraft #92 (@roma-glushko)
- 👷 Publish images to DockerHub #123 (@roma-glushko)
- 🔧 Migrated all API to Fiber #136 (@roma-glushko)
- 👷 Create a image tag with pure version (without distro suffix) #139 (@roma-glushko)

## [0.0.1] (Jan 31st, 2024)

### Added

- ✨Allow to chat message based for specific models #81 (@mkrueger12)

### Changed

- 🔧 Normalize response latency by response token count #78 (@roma-glushko)
- 📝 Added the CLI banner info #112 (@roma-glushko)

### Miscellaneous

- 📝 #114 Make links actual across the project (@roma-glushko)

## [0.0.1-rc.2] (Jan 22nd, 2024)

### Added

- ⚙️ [config] Added validation for config file content #40 (@roma-glushko)
- ⚙️ [config] Allowed to pass HTTP server configs from config file #41 (@roma-glushko)
- 👷 [build] Allowed building Homebrew taps for release candidates #99 (@roma-glushko)

## [0.0.1-rc.1] (Jan 21st, 2024)

### Added
- ✨ [providers] Support for OpenAI Chat API #3 (@mkrueger12)
- ✨ [API] Unified Chat API #54 (@mkrueger12)
- ✨ [providers] Support for Cohere Chat API #5 (@mkrueger12)
- ✨ [providers] Support for Azure OpenAI Chat API #4 (@mkrueger12)
- ✨ [providers] Support for OctoML Chat API #58 (@mkrueger12)
- ✨ [routing] The Routing Mechanism, Adaptive Health Tracking, and Fallbacks #42 #43 #51 (@roma-glushko)
- ✨ [routing] Support for round-robin routing strategy #44 (@roma-glushko)
- ✨ [routing] Support for the least latency routing strategy #46 (@roma-glushko)
- ✨ [routing] Support for weighted round-robin routing strategy #45 (@roma-glushko)
- ✨ [providers] Support for Anthropic Chat API #60 (@mkrueger12)
- ✨ [docs] OpenAPI specifications #22 (@roma-glushko)

### Miscellaneous

- 🔧 [chores] Inited the project #6 (@roma-glushko)
- 🔊 [telemetry] Inited logging  #14 (@roma-glushko)
- 🔧 [chores] Inited Glide's CLI #12 (@roma-glushko)
- 👷 [chores] Setup CI workflows #8 (@roma-glushko)
- ⚙️ [config] Inited configs #11 (@roma-glushko)
- 🔧 [chores] Automatic coverage reports #39 (@roma-glushko)
- 👷 [build] Setup release workflows #9 (@roma-glushko)

[unreleased]: https://github.com/EinStack/glide/compare/0.1.0-rc.1...HEAD
[0.1.0-rc.1]: https://github.com/EinStack/glide/compare/0.0.3...0.1.0-rc.1
[0.0.3]: https://github.com/EinStack/glide/compare/0.0.3-rc.1..0.0.3
[0.0.3-rc.2]: https://github.com/EinStack/glide/compare/0.0.3-rc.1..0.0.3-rc.2
[0.0.3-rc.1]: https://github.com/EinStack/glide/compare/0.0.2..0.0.3-rc.1
[0.0.2]: https://github.com/EinStack/glide/compare/0.0.2-rc.1..0.0.2
[0.0.2-rc.2]: https://github.com/EinStack/glide/compare/0.0.2-rc.1..0.0.2-rc.2
[0.0.2-rc.1]: https://github.com/EinStack/glide/compare/0.0.1..0.0.2-rc.1
[0.0.1]: https://github.com/EinStack/glide/compare/0.0.1-rc.2..0.0.1
[0.0.1-rc.2]: https://github.com/EinStack/glide/compare/0.0.1-rc.1..0.0.1-rc.2
[0.0.1-rc.1]: https://github.com/EinStack/glide/releases/tag/0.0.1-rc.1
