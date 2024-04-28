# Contribution Guide

First off, we are super excited that you are willing to improve Glide 🙌

There are three areas to contribute:

- **Technical**: Help us to improve existing functionality, fix bugs, and bring on new features both to Glide and related repositories like Python SDK.
- **Documentation**: Improve documentation content, uncover undocumented features & gotchas, write guides and walkthroughs. 
- **Vision**: Help us to uncover use cases where Glide could have helped that might be useful for a broader set of people

---

## Technical Contribution

### Communication & Coordination

We value your time. 
To make your onboarding as smooth as possible while reducing amount of back and forth, 
we coordinate and communicate in [the EinStack's Discord space](https://discord.gg/rsBzprY7uT) before jumping on anything major.

Overcommunication is the key to solving many problems.

### GEPs

We are using [enhancement proposals](https://github.com/EinStack/geps) to 
define bigger problems and suggest our solutions to them.

The enhancement proposals share your ideas on solving the problem and let other people give a feedback, 
identify areas to investigate, brainstorm alternatives.

To start a new GEP, you don't have to know all the answers to all questions. 
You can outline gaps and let other people contribute their ideas on possible solutions.

### Dev Commands

Many useful commands are in [the root makefile](Makefile). 
We use make as a convenient interface to automate a bunch of commands like codebase linting, running tests, running dev binary, etc.
Be sure to take a look at all available commands.

### CI Checks

All important checks are automated on the level of pull request checks. 
Be sure to keep your PRs green, before moving the PR to the review stage.

#### build:dry-run

The Glide repository has a special `build:dry-run` label that allows to run the release workflow without actually publishing Glide artefacts. 
This is helpful for:
- testing image building
- making sure any changes to the release workflow works fine

## Improve Our Documentation

### Typos & Uncovered Functionality

If you spot a typo or incorrect information, please do use the `raise issue` or `suggest edits` functionality directly on the documentation page.

If you see some uncovered functionality, please fill briefly [a Github issue](https://github.com/EinStack/docs/issues).

### Guides

A special place takes our guides. Guide is a walkthrough that solves a concrete use case or problem step by step.

To inspire our end users and illustrate the true capabilities of Glide, we want to grow the number of guides. 

If you have any specific use cases to cover, please do let us know in [Discord](https://discord.gg/rsBzprY7uT) or [our docs repo](https://github.com/EinStack/docs) (even if you don't have a chance to work on that).

## Expand Our Vision

If you feel like we have overlooked 
some useful functionality or features that would be great to have, 
feel free to create a new discussion in [our Github Discussions](https://github.com/EinStack/glide/discussions/categories/ideas).

We will review and discuss all ideas and will try to fit them into [the Glide's roadmap](ROADMAP.md).

## Don't want to contribute but uses Glide

That's perfectly fine! 

Feel free to connect with us in [Discord](https://discord.gg/rsBzprY7uT) and ask any question you have.
Remember, there are no dumb questions, but there can be missing opportunities to make your life easier if you don't speak up about things you struggle with.


