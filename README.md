[<img width="70%" src=".github/logo.svg?raw=true"/>](https://letsblock.it)

## Remove low-quality content and useless nags, focus on what matters.

[![CI](https://github.com/letsblockit/letsblockit/actions/workflows/ci.yml/badge.svg)](https://github.com/letsblockit/letsblockit/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/letsblockit/letsblockit)](https://goreportcard.com/report/github.com/letsblockit/letsblockit)
[![Uptime Robot status](https://img.shields.io/uptimerobot/status/m790590578-5dd39cb706475421fade8af1)](https://stats.uptimerobot.com/YxoOltDPJX)

This repository holds the data and server source code for the Let's Block It project. 
Please check out the project's [about page](https://letsblock.it/help/about)
and [contributing guidelines](https://letsblock.it/help/contributing).

## How to use this project

- Go to [letsblock.it](https://letsblock.it), create an account and curate your personal filter list
based on our filter templates. A personal URL will be generated for you to add to your uBlock Origin lists.
- Alternatively, you can maintain your filter parameters in a YAML file, and use [the local render CLI](cmd/render)
to generate your personal filter list. You will still need to host it on a web server for it to be retrievable
by uBlock Origin.
- Self-hosting of [the letsblock.it server](cmd/server) is currently not officially supported, as the database schema is
evolving quickly. Don't hesitate to open an issue to discuss this if you'd be interested.

## Repository structure

Here are the most important folders, check their README.md file for specific documentation:
- **Filter templates** are stored in [data/filters/](data/filters)
- The **local render command** is in [cmd/render/](cmd/render)
- The **letsblock.it server backend** is in [cmd/server/](cmd/server), it depends on golang backend code from [src/](src)
and handlebars page templates from [data/pages/](data/pages)
- The **letsblock.it frontend assets** are defined in [src/assets/](src/assets) and pre-compiled in [data/assets/dist/](data/assets/dist)
- **Development scripts** are in [scripts/](scripts)

## Opening issues

Don't hesitate to open issues on this repository, using the following templates:

- [Improve / fix an existing filter template](https://github.com/letsblockit/letsblockit/issues/new?labels=filter-data&template=update-filter.yaml)
- [Suggest a new filter template](https://github.com/letsblockit/letsblockit/issues/new?labels=filter-data&template=new-filter.yaml)
- [Report a bug](https://github.com/letsblockit/letsblockit/issues/new?labels=&template=bug-report.yaml)
- [Other topics](https://github.com/letsblockit/letsblockit/issues/new)
