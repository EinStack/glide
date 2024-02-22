# Changelog

The changelog consists of three categories:
- **Features** - a new functionality that brings value to users
- **Improvements** - bugfixes, performance and other types of improvements to existing functionality
- **Miscellaneous** - all other updates like build, release, CLI, etc.

## 0.0.2-rc.2, 0.0.2 (Feb 22nd, 2024)

### Features

- ✨ #142: [Lang Chat Router] Ollama Support (@mkrueger12)
- ✨ #131: [Lang Chat Router] AWS Bedrock Support (@mkrueger12)

### Miscellaneous

- 👷 #155 Fixing the dockerhub authorization step in the release workflow (@roma-glushko)
- ♻️  #151: Moved specific provider schemas closer to provider's packages (@roma-glushko)

## 0.0.2-rc.1 (Feb 12th, 2024)

### Features

- ✨#117 Allow to load dotenv files (@roma-glushko)

### Improvements

- ✨👷#91 Support for Windows (@roma-glushko)
- 👷 #139 Build Glide for OpenBSD and ppc65le, s390x, riscv64 architectures (@roma-glushko)

### Miscellaneous

- 👷 #92 Release binaries to Snapcraft (@roma-glushko)
- 👷 #123 publish images to DockerHub (@roma-glushko)
- 🔧 #136 Migrated all API to Fiber (@roma-glushko)
- 👷 #139 Create a image tag with pure version (without distro suffix) (@roma-glushko)

## 0.0.1 (Jan 31st, 2024)

### Features

- ✨ #81: Allow to chat message based for specific models (@mkrueger12)

### Improvements

- 🔧 #78: Normalize response latency by response token count (@roma-glushko)
- 📝 #112 added the CLI banner info (@roma-glushko)

### Miscellaneous

- 📝 #114 Make links actual across the project (@roma-glushko)

## 0.0.1-rc.2 (Jan 22nd, 2024)

### Improvements

- ⚙️ [config] Added validation for config file content #40 (@roma-glushko)
- ⚙️ [config] Allowed to pass HTTP server configs from config file #41 (@roma-glushko)
- 👷 [build] Allowed building Homebrew taps for release candidates #99 (@roma-glushko)

## 0.0.1-rc.1 (Jan 21st, 2024)

### Features
- ✨ [providers] Support for OpenAI Chat API #3 (@mkrueger12)
- ✨ [API] Unified Chat API #54 (@mkrueger12)
- ✨ [providers] Support for Cohere Chat API #5 (@mkrueger12)
- ✨ [providers] Support for Azure OpenAI Chat API #4 (@mkrueger12)
- ✨ [providers] Support for OctoML Chat API #58 (@mkrueger12)
- ✨ [routing] The Routing Mechanism, Adaptive Health Tracking, and Fallbacks #42 #43 #51 (@roma-glushko)
- ✨ [routing] Support for round robin routing strategy #44 (@roma-glushko)
- ✨ [routing] Support for the least latency routing strategy #46 (@roma-glushko)
- ✨ [routing] Support for weighted round robin routing strategy #45 (@roma-glushko)
- ✨ [providers] Support for Anthropic Chat API #60 (@mkrueger12)
- ✨ [docs] OpenAPI specifications #22 (@roma-glushko)

### Miscellaneous

- 🔧 [chores] Inited the project #6 (@roma-glushko)
- 🔊 [telemetry] Inited logging  #14 (@roma-glushko)
- 🔧 [chores] Inited Glide's CLI #12 (@roma-glushko)
- 👷 [chores] Setup CI workflows #8 (@roma-glushko)
- ⚙️ [config] Inited configs #11 (@roma-glushko)
-  🔧 [chores] Automatic coverage reports #39 (@roma-glushko)
- 👷 [build] Setup release workflows #9 (@roma-glushko)
