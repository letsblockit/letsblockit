# letsblock.it frontend sources

This folder holds the npm project to build the frontend assets (CSS and JS).

Because of the low rate of changes, I chose to pre-generate them and store them
in the [data/assets/dist/](../../data/assets/dist) folder. You can use
[the helper script](../../scripts/update-assets.sh) to do so.

## Scope limitations

Please note that the project is designed to work without Javascript enabled, hence the choice of server-side rendering
and [progressive enhancement](https://developer.mozilla.org/en-US/docs/Glossary/Progressive_Enhancement). Frontend
contributions are very welcome (this is not my area of expertise), but need to stay in this spirit.
