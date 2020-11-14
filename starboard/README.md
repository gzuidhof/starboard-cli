# starboard (cli)

To download a new runtime, specify its version and the target folder in the `go:generate` in main.go:

```bash
go run scripts/download_runtime/main.go 0.6.4 web/static/vendor/
```


## Development
To run the serve command with the latest static assets and templates without having to `go generate`, use:

```
go run main.go serve --static_folder web/static --templates_folder web/templates
```
