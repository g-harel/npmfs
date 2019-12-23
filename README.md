<!--

todo
    add download link to all dirs/files
    github badge?
    registry picker?

-->

<!-- REMEMBER TO MIRROR CONTENT CHANGES ON HOMEPAGE -->

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

<details>
    <summary>Examples</summary>
    <a href="https://npmfs.com/lodash">https://npmfs.com/lodash</a>
    <br>
    <a href="https://npmfs.com/package/lodash">https://npmfs.com/package/lodash</a>
    <br>
    <a href="https://npmfs.com/request">https://npmfs.com/request</a>
    <br>
    <a href="https://npmfs.com/package/request">https://npmfs.com/package/request</a>
</details>

##

Specific package versions can be accessed directly.

```
https://npmfs.com/package/<name>/<version>
```

<details>
    <summary>Examples</summary>
    <a href="https://npmfs.com/package/chalk/2.4.2">https://npmfs.com/package/chalk/2.4.2</a>
    <br>
    <a href="https://npmfs.com/package/chalk/1.0.0">https://npmfs.com/package/chalk/1.0.0</a>
    <br>
    <a href="https://npmfs.com/package/commander/2.20.0">https://npmfs.com/package/commander/2.20.0</a>
    <br>
    <a href="https://npmfs.com/package/commander/1.0.0">https://npmfs.com/package/commander/1.0.0</a>
</details>

##

Directories and files inside the package are viewed by appending the path.

```
https://npmfs.com/package/<name>/<version>/<path>
https://npmfs.com/package/<name>/<version>/example/
https://npmfs.com/package/<name>/<version>/example/index.js
```

<details>
    <summary>Examples</summary>
    <a href="https://npmfs.com/package/async/3.0.1/internal/">https://npmfs.com/package/async/3.0.1/internal/</a>
    <br>
    <a href="https://npmfs.com/package/async/3.0.1/internal/once.js">https://npmfs.com/package/async/3.0.1/internal/once.js</a>
    <br>
    <a href="https://npmfs.com/package/react/16.8.6/umd/">https://npmfs.com/package/react/16.8.6/umd/</a>
    <br>
    <a href="https://npmfs.com/package/react/16.8.6/package.json">https://npmfs.com/package/react/16.8.6/package.json</a>
</details>

##

Package versions are compared (`version-0 .. version-1`) by navigating to the root directory of `version-0`, clicking on the `diff` link in the top right, and selecting `version-1` in the version list.

```
https://npmfs.com/compare/<name>/<version-0>/<version-1>
```

In this compare view, line numbers are a shortcut to their respective file's source.

<details>
    <summary>Examples</summary>
    <a href="https://npmfs.com/compare/debug/4.1.0/4.1.1">https://npmfs.com/compare/debug/4.1.0/4.1.1</a>
    <br>
    <a href="https://npmfs.com/compare/debug/3.0.0/4.0.0">https://npmfs.com/compare/debug/3.0.0/4.0.0</a>
    <br>
    <a href="https://npmfs.com/compare/underscore/1.9.0/1.9.1">https://npmfs.com/compare/underscore/1.9.0/1.9.1</a>
    <br>
    <a href="https://npmfs.com/compare/underscore/1.6.0/1.7.0">https://npmfs.com/compare/underscore/1.6.0/1.7.0</a>
</details>

##

Deep links are found on the right side of lines in both the file and diff views. These links will add a hash to the url which scrolls the browser to the selected line, and highlights it.

```
https://npmfs.com/package/<name>/<version>/index.js#<line>
https://npmfs.com/compare/<name>/<version-0>/<version-1>#<line>
```

<details>
    <summary>Examples</summary>
    <a href="https://npmfs.com/package/bluebird/3.5.5/js/release/race.js#L32">https://npmfs.com/package/bluebird/3.5.5/js/release/race.js#L32</a>
    <br>
    <a href="https://npmfs.com/compare/bluebird/3.5.4/3.5.5#D0L35">https://npmfs.com/compare/bluebird/3.5.4/3.5.5#D0L35</a>
    <br>
    <a href="https://npmfs.com/package/moment/2.24.0/locale/ca.js#L81">https://npmfs.com/package/moment/2.24.0/locale/ca.js#L81</a>
    <br>
    <a href="https://npmfs.com/compare/moment/2.23.0/2.24.0#D19L0">https://npmfs.com/compare/moment/2.23.0/2.24.0#D19L0</a>
</details>

##

The sticky header sections often contain useful links to navigate between pages.

## Development

```bash
$ go run main.go
```

_Requires a least `go1.11` (for go modules support) and `git` to be installed._

_Replace `go` with [`gin`](https://github.com/codegangsta/gin) for auto-restarts._

##

Tests are run using the standard go test command (`go test ./...`).

## Deployment

This project is hosted on [Google Cloud Platform](https://cloud.google.com/)'s [Cloud Run](https://cloud.google.com/run/) product _(currently in beta)_.

The [deployment workflow](./.github/main.workflow) uses [GitHub Actions](https://developer.github.com/actions/) to publish a new image and update the running service on push.

## License

[MIT](./LICENSE)
