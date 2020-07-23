<!--

todo
    github badge?
    registry picker?

-->

# [npmfs](https://npmfs.com)

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
