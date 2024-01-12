# Glide: Cloud-Native LLM Gateway for Seamless LLMLOps
<div align="center">
    <img src="docs/images/glide.png" width="400px" alt="Glide GH Header" />
</div>

[![LICENSE](https://img.shields.io/github/license/modelgateway/glide.svg?style=flat-square&color=%233f90c8)](https://github.com/modelgateway/glide/blob/main/LICENSE)
[![codecov](https://codecov.io/github/modelgateway/glide/graph/badge.svg?token=F7JT39RHX9)](https://codecov.io/github/modelgateway/glide)

Glide is your go-to cloud-native model gateway, delivering high-performance LLMLOps in a lightweight, all-in-one package.

We take all problems of managing and communicating with external providers out of your applications,
so you can dive into tackling your core challenges.

Glide sits between your application and model providers to seamlessly handle various LLMOps tasks like
model failover, caching, key management, etc. 

> [!Warning]
> Glide is under active development right now. Give us a star to support the project ✨

## Features

- **Unified REST API** across providers. Avoid vendor lock-in and changes in your applications when you swap model providers.
- **High availability** and **resiliency** when working with external model providers. Automatic **fallbacks** on provider failures, rate limits, transient errors. Smart retries to reduce communication latency.
- Support **popular LLM providers**.
- **High performance**. Performance is our priority. We want to keep Glide "invisible" for your latency-wise, while providing rich functionality.
- **Production-ready observability** via OpenTelemetry, emit metrics on models health, allows whitebox monitoring.
- Straightforward and simple maintenance and configuration, centrilized API key control & management & rotation, etc.

## Supported Providers

### Large Language Models

|                                                     | Provider      | Support Status  |
|-----------------------------------------------------|---------------|-----------------|
| <img src="docs/images/openai.svg" width="18" />     | OpenAI        | 👍  Supported  |
| <img src="docs/images/azure.svg" width="18" />      | Azure OpenAI  | 🏗️ Coming Soon  |
| <img src="docs/images/cohere.png" width="18" />     | Cohere        | 🏗️ Coming Soon  |
| <img src="docs/images/octo.png" width="18" />     | OctoML        | 🏗️ Coming Soon  |
| <img src="docs/images/anthropic.svg" width="18" />  | Anthropic     | 🏗️ Coming Soon |
| <img src="docs/images/bard.svg" width="18" />       | Google Gemini | 🏗️ Coming Soon |


## Get Started

TBU

### API Docs

Glide comes with OpenAPI documentation that is accessible via http://127.0.0.1:9099/v1/swagger/index.html

## Community

- Join [Discord](https://discord.gg/z4DmAbJP) for real-time discussion

Open [an issue](https://github.com/modelgateway/glide/issues) or start [a discussion](https://github.com/modelgateway/glide/discussions) 
if there is a feature or an enhancement you'd like to see in Glide.

## Contribute

- Maintainers
    
    - [Roman Hlushko](https://github.com/roma-glushko), Software Engineer, Distributed Systems & MLOps
    - [Max Krueger](https://github.com/mkrueger12), Data & ML Engineer

Thanks everyone for already put their effort to make Glide better and more feature-rich: 

<a href="https://github.com/modelgateway/glide/graphs/contributors">
  <img src="https://contributors-img.web.app/image?repo=modelgateway/glide" />
</a>
