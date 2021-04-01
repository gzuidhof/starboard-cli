# starboard (cli)

To download a new runtime, specify its version and the target folder in the `go:generate` in main.go.

## Development
To run the serve command with the latest static assets and templates without having to `go generate`, use:

```
go run main.go serve --static_folder web/static --templates_folder web/templates
```

Consider it live-reload as long as you are only changing template files :).

## Releases

Releases are minted on CI, you can create one locally by running
```
goreleaser --snapshot --skip-publish --rm-dist
```
