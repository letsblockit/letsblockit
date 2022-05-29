# Filter list rendering CLI

As an alternative to using [the hosted Let's Block It service](https://letsblock.it),
you can use this CLI to render your list to a file, sourcing the settings from a YAML file.

You can export your current list from [you account page](https://letsblock.it/user/account) or start
from an empty file.

## Example input file

```yaml
title: "example list"               # Title is mandatory
instances:                          # List of instances, multiple instances of a given filter is supported
  - filter: google-search-cleanup   # Filter name, matching the template's file name
    params:                         # Map of param name -> param value
      only-results: false
      related-questions: true
      related-searches: true
  - filter: youtube-shorts          # Params is optional if the filter has no options
```

Check the available filters and parameters by browsing
[the filter sources](https://github.com/letsblockit/letsblockit/tree/main/data/filters).

## How to run the render

[Install the golang compiler](https://go.dev/doc/install), then run:

```shell
go run github.com/letsblockit/letsblockit/cmd/render@latest my-list.yaml > output.txt
```
