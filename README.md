[<img width="70%" src=".github/logo.svg?raw=true"/>](https://letsblock.it)

## Remove low-quality content and useless nags, focus on what matters.

[![CI](https://github.com/letsblockit/letsblockit/actions/workflows/ci.yml/badge.svg)](https://github.com/letsblockit/letsblockit/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/letsblockit/letsblockit)](https://goreportcard.com/report/github.com/letsblockit/letsblockit)
[![Uptime Robot status](https://img.shields.io/uptimerobot/status/m790590578-5dd39cb706475421fade8af1)](https://stats.uptimerobot.com/YxoOltDPJX)

This repository holds the data and server source code for the Let's Block It project. 
Please check out the project's [about page](https://letsblock.it/help/about)
and [contributing guidelines](https://letsblock.it/help/contributing).

## How to use this project

There are three ways to use our filter templates:

1. Go to [letsblock.it](https://letsblock.it), create an account and curate your personal filter list
based on our filter templates. A personal URL will be generated for you to add to your uBlock Origin lists.
2. Self-hosting of [the letsblock.it server](cmd/server) is now supported for experienced users, read
  [the setup doc](cmd/server/README.md) to get started.
3. Alternatively, you can maintain your filter parameters in a YAML file, and use [the local render CLI](cmd/render)
to generate your personal filter list. You will still need to host it on a web server for it to be retrievable
by uBlock Origin.

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

## Thanks to our contributors

Thanks goes to these wonderful people ([emoji key](https://allcontributors.org/docs/en/emoji-key)):

<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->
<!-- prettier-ignore-start -->
<!-- markdownlint-disable -->
<table>
  <tbody>
    <tr>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/xvello"><img src="https://avatars.githubusercontent.com/u/6241083?v=4?s=100" width="100px;" alt="Xavier Vello"/><br /><sub><b>Xavier Vello</b></sub></a><br /><a href="#maintenance-xvello" title="Maintenance">🚧</a> <a href="#data-xvello" title="Data">🔣</a> <a href="https://github.com/letsblockit/letsblockit/commits?author=xvello" title="Code">💻</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/JohnyP36"><img src="https://avatars.githubusercontent.com/u/81161435?v=4?s=100" width="100px;" alt="JohnyP36"/><br /><sub><b>JohnyP36</b></sub></a><br /><a href="#maintenance-JohnyP36" title="Maintenance">🚧</a> <a href="#data-JohnyP36" title="Data">🔣</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/BPower0036"><img src="https://avatars.githubusercontent.com/u/80090789?v=4?s=100" width="100px;" alt="BPower0036"/><br /><sub><b>BPower0036</b></sub></a><br /><a href="#maintenance-BPower0036" title="Maintenance">🚧</a> <a href="#data-BPower0036" title="Data">🔣</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/quenhus"><img src="https://avatars.githubusercontent.com/u/8274288?v=4?s=100" width="100px;" alt="quenhus"/><br /><sub><b>quenhus</b></sub></a><br /><a href="#data-quenhus" title="Data">🔣</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/androidkotlindev"><img src="https://avatars.githubusercontent.com/u/67307199?v=4?s=100" width="100px;" alt="neonItem"/><br /><sub><b>neonItem</b></sub></a><br /><a href="#data-androidkotlindev" title="Data">🔣</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/xbc5"><img src="https://avatars.githubusercontent.com/u/87829033?v=4?s=100" width="100px;" alt="xbc5"/><br /><sub><b>xbc5</b></sub></a><br /><a href="#data-xbc5" title="Data">🔣</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/Markussss"><img src="https://avatars.githubusercontent.com/u/43534?v=4?s=100" width="100px;" alt="Markus Igeland"/><br /><sub><b>Markus Igeland</b></sub></a><br /><a href="#data-Markussss" title="Data">🔣</a></td>
    </tr>
    <tr>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/kiriya-aoi"><img src="https://avatars.githubusercontent.com/u/34067059?v=4?s=100" width="100px;" alt="kiriya-aoi"/><br /><sub><b>kiriya-aoi</b></sub></a><br /><a href="#data-kiriya-aoi" title="Data">🔣</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/krisu5"><img src="https://avatars.githubusercontent.com/u/3540275?v=4?s=100" width="100px;" alt="Kristian Salonen"/><br /><sub><b>Kristian Salonen</b></sub></a><br /><a href="#data-krisu5" title="Data">🔣</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/mailyspamy"><img src="https://avatars.githubusercontent.com/u/25592998?v=4?s=100" width="100px;" alt="mailyspamy"/><br /><sub><b>mailyspamy</b></sub></a><br /><a href="#data-mailyspamy" title="Data">🔣</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/santilococo"><img src="https://avatars.githubusercontent.com/u/27114158?v=4?s=100" width="100px;" alt="Santiago Lo Coco"/><br /><sub><b>Santiago Lo Coco</b></sub></a><br /><a href="#data-santilococo" title="Data">🔣</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://opencollective.com/letsblockit#section-contributors"><img src="https://images.opencollective.com/guest-1bcb213b/avatar.png?s=100" width="100px;" alt="Webswonder"/><br /><sub><b>Webswonder</b></sub></a><br /><a href="#financial-Webswonder" title="Financial">💵</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/SiggyPony"><img src="https://avatars.githubusercontent.com/u/11322269?v=4?s=100" width="100px;" alt="Sigourney Juneau"/><br /><sub><b>Sigourney Juneau</b></sub></a><br /><a href="#data-SiggyPony" title="Data">🔣</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://github.com/alexanderadam"><img src="https://avatars.githubusercontent.com/u/372620?v=4?s=100" width="100px;" alt="Alexander ADAM"/><br /><sub><b>Alexander ADAM</b></sub></a><br /><a href="#promotion-alexanderadam" title="Promotion">📣</a></td>
    </tr>
  </tbody>
</table>

<!-- markdownlint-restore -->
<!-- prettier-ignore-end -->

<!-- ALL-CONTRIBUTORS-LIST:END -->

## Thanks to our sponsors

<table>
  <tr>
    <td align="center"><a href="https://www.ory.sh/"><img src=".github/logo-ory.svg?raw=true"/></a></td>
    <td>Thanks to Ory for hosting the official instance's<br/>
        user database on <a href="https://www.ory.sh/cloud/">Ory Cloud.</a></td>
  </tr>
  <tr>
    <td align="center"><a href="https://www.jetbrains.com/opensource/"><img src=".github/logo-jetbrains.svg?raw=true"/></a></td>
    <td>Thanks to Jetbrains for providing free GoLand licenses.</td>
  </tr>
</table>
