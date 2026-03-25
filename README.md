# stockuploader

stockuploader allows you to upload images and video via command line to Adobe Stock. Simple and easy. While the stockuploader works, it is a very early version.

## Usage

```
stockuploader

        Usage:
          stockuploader [options] <file>...
          stockuploader (-v | --version)
          stockuploader (-h | --help)

        Arguments:
          <file>                      One or multiple files you'd like to upload

        Options:
          -h --help                   Show this screen.
          -v --version                Show version number.
          -c --config=<path>         Load username, password, and optional remote/port from TOML (overrides auto-detection).
          -u --username=<username>  Provide your Adobe username (overrides config).
          -p --password=<password>  Provide your Adobe password (overrides config).
          -r --remote=<host>        SFTP remote host (when omitted: config value, else Adobe Stock SFTP host).
          -P --port=<port>          Remote port (when omitted: config value, else 22).
```

### TOML config

You can provide credentials via a TOML file with `--config=<path>`, or (if `--config` is omitted) stockuploader will auto-detect the first existing config file in this order:

1. `./stockuploader.toml`
2. `~/.config/stockuploader.toml`
3. `~/.stockuploader.toml`
4. `/etc/stockuploader.toml`

The config supports these keys (as strings):

- `username`
- `password`
- `remote` (optional)
- `port` (optional)

`--username` / `--password` override the values from the config; `--remote` / `--port` override `remote` / `port`.

## Installation

If you are on a Mac and you are using homebrew, then you can install stockuploader via my tap.

`$ brew tap brejoc/tap`  
`$ brew install stockuploader`

## Building

Install the most recent version of Go. Currently this is version `1.21.1`. Since all of the dependencies are vendorized, a `go build` is enough. After that you'll see a `stockuploader` binary in the project folder.