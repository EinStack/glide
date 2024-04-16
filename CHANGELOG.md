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

## [0.0.3-rc2], [0.0.3] (Apr 17, 2024)

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

[unreleased]: https://github.com/olivierlacan/keep-a-changelog/compare/0.0.3...HEAD
[0.0.3]: https://github.com/EinStack/glide/compare/0.0.3-rc.1..0.0.3
[0.0.3-rc.2]: https://github.com/EinStack/glide/compare/0.0.3-rc.1..0.0.3-rc.2
[0.0.3-rc.1]: https://github.com/EinStack/glide/compare/0.0.2..0.0.3-rc.1
[0.0.2]: https://github.com/EinStack/glide/compare/0.0.2-rc.1..0.0.2
[0.0.2-rc.2]: https://github.com/EinStack/glide/compare/0.0.2-rc.1..0.0.2-rc.2
[0.0.2-rc.1]: https://github.com/EinStack/glide/compare/0.0.1..0.0.2-rc.1
[0.0.1]: https://github.com/EinStack/glide/compare/0.0.1-rc.2..0.0.1
[0.0.1-rc.2]: https://github.com/EinStack/glide/compare/0.0.1-rc.1..0.0.1-rc.2
[0.0.1-rc.1]: https://github.com/EinStack/glide/releases/tag/0.0.1-rc.1
