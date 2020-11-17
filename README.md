starboard-cli
---

The `starboard` command line tool is used for interacting with [Starboard Notebooks](https://github.com/gzuidhof/starboard-notebook) locally. It starts a local webserver so you can view and edit notebook files on your computer.

## Installation
```bash
go get github.com/gzuidhof/starboard-cli/starboard
```

Pre-built binaries and a NPM distribution will be available later.

## Usage
```bash
# Serve files in current folder
starboard serve .

## Serve files under a certain path
starboard serve ./assets
```

## Screenshots

### `starboard serve` directory listing
![Starboard Serve Screenshot](https://i.imgur.com/6k8VDz8.png)

### `starboard serve` notebook view/editor
![Starboard Server Notebook View Screenshot](https://i.imgur.com/gy2Iuyl.png)

## License
This is free software; you can redistribute it and/or modify it under the terms of the [Mozilla Public License Version 2.0](./LICENSE).
