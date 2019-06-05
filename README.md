<!--

todo
    code documentation
    write more tests
    put behind cdn
    track popularity
    registry picker

The linked repository can have entirely different contents (This discrepancy is often for good reason, like removing unnecessary files or publishing bundled, transpiled, or minified files).

-->

# [npmfs](https://npmfs.com)

JavaScript Package Inspector

- [**Package Viewer** – Browse package files at any version.](https://npmfs.com/package/express/4.17.1/)

- [**Package Diff** – Compare packages across versions.](https://npmfs.com/compare/express/4.17.0/4.17.1/)

- [**Deep Links** – Link to specific lines in files and diffs.](https://npmfs.com/package/express/4.17.1/Readme.md#L16)

## Motivations

Published packages don't always include a link to their source on GitHub.

The linked repository is not necessarily representative of published package.

## Usage

Package root provides links to all available versions, and highlights the one tagged as `latest`.

```
https://npmfs.com/<name>
https://npmfs.com/package/<name>
```

_When viewing a package on npm, it is conveniently accessible with a one-char edit to the page url._

```diff
- https://www.npmjs.com/package/<name>
+ https://www.npmfs.com/package/<name>
```

##

Specific package versions can be accessed directly.

```
https://npmfs.com/package/<name>/<version>
```

##

Directories and files inside the package are viewed by appending the path.

```
https://npmfs.com/package/<name>/<version>/<path>
https://npmfs.com/package/<name>/<version>/example/
https://npmfs.com/package/<name>/<version>/example/index.js
```

##

Package versions are compared (`version-0 .. version-1`) by navigating to the root directory of `version-0`, clicking on the `diff` link in the top right, and selecting `version-1` in the version list.

```
https://npmfs.com/compare/<name>/<version-0>/<version-1>
```

In this compare view, line numbers are a shortcut to their respective file's source.

##

Deep links are found on the right side of lines in both the file and diff views. These links will add a hash to the url which scrolls the browser to the selected line, and highlights it.

```
https://npmfs.com/package/<name>/<version>/index.js#<line>
https://npmfs.com/compare/<name>/<version-0>/<version-1>#<line>
```

##

The sticky header sections often contain useful links to navigate between pages.

## Development

```bash
$ go run main.go
```

_Requires a least `go1.11` (for go modules support) and `git` to be installed._

##

Tests are run using the standard go test command (`go test ./...`).

## Deployment

This project is hosted on [Google Cloud Platform](https://cloud.google.com/)'s [Cloud Run](https://cloud.google.com/run/) product _(currently in beta)_.

The [deployment workflow](./.github/main.workflow) uses [GitHub Actions](https://developer.github.com/actions/) to publish a new image and update the running service on push.

## License

[MIT](./LICENSE)
