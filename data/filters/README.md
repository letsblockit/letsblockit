# Filter templates

This folder holds the filter templates for [https://letsblock.it](https://letsblock.it), one filter per file. Adding one
filter is as easy as adding a new file and opening a PR. Please read
the [contribution guidelines](https://letsblock.it/help/contributing) before getting started.

## File structure

Each file is split into two sections, separated by a `---` line:

- the first half is the filter properties, as a `YAML` object, described in the next section,
- the second half is the filter description, as [Markdown](https://www.markdownguide.org/basic-syntax). It will be
  displayed above the filter parameters form.

## YAML properties

Let's start with the easy properties first: `title` and `tags` (a list) are self-explanatory, they
are used in the filter list and the filter view pages.

Then comes the `params` list. It describes what parameters the filter accepts. Each entry must have:

- A lowercase `name`, used in the template rendering
- A human-facing `description`, used in the parameters form
- A `type`: either `checkbox` (true/false boolean), `string` (a single input line), `list` (several input lines, user
  can add / remove more dynamically) or `multiline` (a multi-line input box, used in `custom-rules` for example)
- A `default` value, used when displaying the filter to a new user. Valid values are:
    - `checkbox`: `true` or `false`
    - `list` and `multiline`: a single string between double quotes
    - `list`: a list of strings, in the YAML list format (e.g. `["one", "two"]`)

There are two ways to define what content blocking rules to output:

- For simple templates with only checkbox parameters, add a `rules` field to each parameter, for rules
  to add when that parameter is enabled

- If other parameter types are used, you'll pass a `template` property, in [Handlebars](https://handlebarsjs.com/guide/)
  format. Every param is accessible via their `name`. To check that your syntax is correct and guard against regression,
  `tests` cases are written as a list of objects, with the following fields:
    - `params` is a key/value object, defining the test case input
    - `output` is the expected output of the filter template, the test will fail if it differs

If you have the Go compiler [installed](https://go.dev/doc/install), you can run `go test -v ./src/filters/`
in the project's root directory. The tests will validate the filters' format and syntax, and run their test cases.
Otherwise, they will run on your PR when it is reviewed.

## uBlock Origin filter syntax

Have a look at the following pages to learn about the syntax:

- The uBlock Origin wiki pages: [Static filter syntax](https://github.com/gorhill/uBlock/wiki/Static-filter-syntax) and
  [Procedural cosmetic filters](https://github.com/gorhill/uBlock/wiki/Procedural-cosmetic-filters).
- [This cheatsheet](https://github.com/DandelionSprout/adfilt/blob/master/Wiki/SyntaxMeaningsThatAreActuallyHumanReadable.md)
  by DandelionSprout.

## Testing your changes

Filter files are checked by automated tests that will run on your PR once approved.
You can run them locally with one of these commands:

- if you have [the golang compiler](https://go.dev/doc/install) installed: `go test ./src/filters`
- if you have docker installed: `docker run --rm -w /source -v $PWD:/source golang go test ./src/filters`
