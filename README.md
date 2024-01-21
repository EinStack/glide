# Glide: Cloud-Native LLM Gateway for Seamless LLMOps
<div align="center">
    <img src="docs/images/glide.png" width="400px" alt="Glide GH Header" />
</div>

[![LICENSE](https://img.shields.io/github/license/modelgateway/glide.svg?style=flat-square&color=%233f90c8)](https://github.com/modelgateway/glide/blob/main/LICENSE)
[![codecov](https://codecov.io/github/EinStack/glide/graph/badge.svg?token=F7JT39RHX9)](https://codecov.io/github/EinStack/glide)

Glide is your go-to cloud-native LLM gateway, delivering high-performance LLMOps in a lightweight, all-in-one package.

We take all problems of managing and communicating with external providers out of your applications,
so you can dive into tackling your core challenges.

Glide sits between your application and model providers to seamlessly handle various LLMOps tasks like
model failover, caching, key management, etc. 

Take a look at the develop branch.

Check out our [documentation](https://backlandlabs.mintlify.app/introduction)!

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
| <img src="docs/images/azure.svg" width="18" />      | Azure OpenAI  | 👍  Supported  |
| <img src="docs/images/cohere.png" width="18" />     | Cohere        | 👍  Supported |
| <img src="docs/images/octo.png" width="18" />     | OctoML        | 👍  Supported  |
| <img src="docs/images/anthropic.svg" width="18" />  | Anthropic     | 👍  Supported  |
| <img src="docs/images/bard.svg" width="18" />       | Google Gemini | 🏗️ Coming Soon |


### Routers

Routers are a core functionality of Glide. Think of routers as a group of models with some predefined logic. For example, the resilience router allows a user to define a set of backup models should the initial model fail. Another example, would be to leverage the least-latency router to make latency sensitive LLM calls in the most efficient manner.

Detailed info on routers can be found [here](https://backlandlabs.mintlify.app/essentials/routers).

#### Available Routers

| Router      | Description  |
|---------------|-----------------|
| Priority        | When the target model fails the request is sent to the secondary model. The entire service instance keeps track of the number of failures for a specific model reducing latency upon model failure  |
| Least Latency        | This router selects the model with the lowest average latency over time. If the least latency model becomes unhealthy, it will pick the second the best, etc.  |
| Round Robin        | Split traffic equally among specified models. Great for A/B testing.  |
| Weighted Round Robin | Split traffic based on weights. For example, 70% of traffic to Model A and 30% of traffic to Model B.  |


## Get Started

#### Install

The easiest way to deploy Glide is to build from source.

Steps to build a container with Docker can be found [here](https://backlandlabs.mintlify.app/introduction#install-and-deploy).

#### Set Configuration File

Find detailed information on configuration [here](https://backlandlabs.mintlify.app/essentials/configuration).

```yaml
telemetry:
  logging:
    level: debug  # debug, info, warn, error, fatal
    encoding: console

routers:
  language:
    - id: myrouter
      models:
        - id: openai
          openai:
            api_key: ""
```

#### Sample API Request to `/chat` endpoint

See [API Reference](https://backlandlabs.mintlify.app/api-reference/introduction) for more details.

```json
{
 "message":
      {
        "role": "user",
        "content": "Where was it played?"
      },
    "messageHistory": [
      {"role": "system", "content": "You are a helpful assistant."},
      {"role": "user", "content": "Who won the world series in 2020?"},
      {"role": "assistant", "content": "The Los Angeles Dodgers won the World Series in 2020."}
    ]
}
```

### API Docs

Once deployed, Glide comes with OpenAPI documentation that is accessible via http://127.0.0.1:9099/v1/swagger/index.html

---

Other ways to install Glide are available:

### Homebrew (MacOS)

Coming Soon

### Snapcraft (Linux)

Coming Soon

### Docker Images

Glide provides official images in our [GHCR](https://github.com/EinStack/glide/pkgs/container/glide):

- Alpine 3.19:
```bash
docker pull ghcr.io/einstack/glide:latest-alpine 
```

- Ubuntu 22.04 LTS:
```bash
docker pull ghcr.io/einstack/glide:latest-ubuntu
```

- Google Distroless (non-root)
```bash
docker pull ghcr.io/einstack/glide:latest-distroless
```

- RedHat UBI 8.9 Micro
```bash
docker pull ghcr.io/einstack/glide:latest-redhat
```

### Helm Chart

Coming Soon

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
