## Contributing to the project

Thanks for your interest in moving this project forward! If you have any question or suggestion, don't hesitate to
[open an issue on GitHub](https://github.com/letsblockit/letsblockit/issues) for any question or bug, or send us an email
at `hello@<this-domain>.it`.

### Financial contributions

To help keep the official instance running, you can sponsor the project on
[OpenCollective](https://opencollective.com/letsblockit) or
[Github Sponsors](https://github.com/sponsors/letsblockit).

### Contributing filters {#filters}

The best way you can help the project is by contributing new filters or improving existing ones:

- If you know how to open a pull request, check out
  [the syntax documentation](https://github.com/letsblockit/letsblockit/blob/main/data/filters/README.md) and
  [the existing filter sources](https://github.com/letsblockit/letsblockit/tree/main/data/filters). Please enable the
  `Allow edits by maintainers` option on your PR to allow me to help fix syntax issues.
- Don't hesitate to open a GitHub issue to suggest filter improvements, [open an account](https://github.com/join)
  and [use the relevant issue template](https://github.com/letsblockit/letsblockit/issues/new/choose).

*Please note the following scope limitations:*

- blocking ads and trackers is out of the scope of these filters, please contribute to
  [EasyList](https://easylist.to/) or [AdGuard](https://kb.adguard.com/en/general/adguard-ad-filters) instead
- support for content blockers other than `uBlock Origin` is left out for now, to limit the maintenance cost
- as a protest against the toxic endeavours of the Facebook / Meta company, no filter improving any of its websites will
  be accepted in the project.

### Helping with the code {#code}

A detailed contribution documentation will come soon, but in the meantime, you can take a look at the repository:

- the page sources are located in `data/pages`, in [Handlebars](https://handlebarsjs.com/guide/) and Markdown formats
- the golang backend code is in the `src/` folders
- if you want to run the server locally, you will need to apply the postgres schema at `src/db/schema.sql` first
- the frontend sources are in `src/assets`, assets are pre-generated with npm into `data/assets`

Please note that the website is designed to work without Javascript enabled, hence the choice of server-side rendering
and [progressive enhancement](https://developer.mozilla.org/en-US/docs/Glossary/Progressive_Enhancement). Frontend
contributions are very welcome (this is not my area of expertise), but need to stay in this spirit.

### Licencing of your contributions

Both the code and data for this project are licensed under the
[Apache License version 2.0](https://github.com/letsblockit/letsblockit/blob/main/LICENSE.txt), which contains the following
clause:

<div class="alert alert-secondary me-md-5 ms-md-5" role="alert">
Unless You explicitly state otherwise, any Contribution intentionally submitted for inclusion in the Work by You
to the Licensor shall be under the terms and conditions of this License, without any additional terms or conditions.
</div>

This means that by opening a PR or an issue on the project, you accept the terms of this license and allow reuse of your
work under these terms.
