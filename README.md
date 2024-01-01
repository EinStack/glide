# Glide
<div align="center">
    <img src="docs/images/glide.png" width="400px" alt="Glide GH Header" />
</div>

[![LICENSE](https://img.shields.io/github/license/modelgateway/glide.svg?style=flat-square&color=%233f90c8)](https://github.com/modelgateway/glide/blob/main/LICENSE)
[![codecov](https://codecov.io/github/modelgateway/glide/graph/badge.svg?token=F7JT39RHX9)](https://codecov.io/github/modelgateway/glide)

Glide is a cloud-native open source high-performant model gateway. All LLMOps you needed is packed in one lightweight service.

We take all problems and toll of managing and communicating with external providers out of your applications,
so you could focus solving your core problems.

Glide sits between your application and model providers that you use to seamlessly handle various LLMOps tasks like
model failover, caching, etc. 

> [!Warning]
> Glide is under active development right now. Give us a star to support the project ✨

## Features

- **Unified REST API** across providers. Avoid vendor lock-ins and changes in your applications when you adopt new providers.
- **High availability** and **resiliency** working with external model providers. Automatic **fallbacks** on provider failures, rate limits, transient errors. Smart retries to reduce communication latency.
- Support **popular LLM providers**.
- **High performance**. Performance is our priority. We want to keep Glide "invisible" for your latency-wise, while providing rich functionality.
- **Production-ready observability** via OpenTelemetry, emit metrics on models health, allows whitebox monitoring.
- Straightforward and simple maintenance and configuration, API key rotation, etc.

## Supported Providers

### Large Language Models

|                                                     | Provider      | Support Status  |
|-----------------------------------------------------|---------------|-----------------|
| <img src="docs/images/openai.svg" width="18" />     | OpenAI        | 🏗️ Coming Soon |
| <img src="docs/images/azure.svg" width="18" />      | Azure OpenAI  | 🏗️ Coming Soon |
| <img src="docs/images/anthropic.svg" width="18" />  | Anthropic     | 🏗️ Coming Soon |
| <img src="docs/images/cohere.png" width="18" />     | Cohere        | 🏗️ Coming Soon |
| <img src="docs/images/bard.svg" width="18" />       | Google Gemini | 🏗️ Coming Soon |
| <img src="docs/images/localai.webp" width="18" />   | LocalAI       | 🏗️ Coming Soon |


## Get Started

TBU

### API Docs

Glide comes with OpenAPI documentation that could be accessible via http://127.0.0.1:9099/v1/swagger/index.html

## Roadmap

### MVP (Coming soon)

- Unified LLM Chat REST API
- Support for most popular LLM providers
- Seamless model fallbacking
- The Main Load Balancing: Priority, Round Robin, Weighted Round Robin, Latency

### Future

- Exact & Semantic Caching
- Cost Management & Budgeting
- and many more!

Open [an issue](https://github.com/modelgateway/glide/issues) or start [a discussion](https://github.com/modelgateway/glide/discussions) 
if there is a feature or an enhancement you'd like to see in Glide.

## Community

- Join Discord for real-time discussions

## Contribute

- Maintainers
    
    - [Roman Hlushko](https://github.com/roma-glushko), Software Engineer, Distributed Systems & MLOps
    - [Max Krueger](https://github.com/mkrueger12), Data Engineer, Data Scientist

Thanks everyone for already put their effort to make Glide better and more feature-rich: 

<a href="https://github.com/modelgateway/glide/graphs/contributors">
  <img src="https://contributors-img.web.app/image?repo=modelgateway/glide" />
</a>

