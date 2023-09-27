# stockuploader

stockuploader allows you to upload images and video via command line to Adobe Stock. Simple and easy. While the stockuploader works, it is a very early version.

## Usage

```
stockuploader

        Usage:
          stockuploader --username peter --password secret <file>...
          stockuploader --username peter --password secret [--remote hostname] <file>...
          stockuploader --username peter --password secret [--remote hostname] [--port 22] <file>...
          stockuploader (-v | --version)

          stockuploader (-h | --help)

        Arguments:
          <file>                One or multiple files you'd like to upload

        Options:
          -h --help                   Show this screen.
          -v --version                Show version number.
          -u --username=<username>    Provide your Adobe username.
          -p --password=<password>    Provide your Adobe password.
          -r --remote=<host>          sFTP remote host [default: sftp.contributor.adobestock.com].
          -P --port=<port>            Remote port [default: 22].
```

## Installation

If you are on a Mac and you are using homebrew, then you can install stockuploader via my tap.

`$ brew tap brejoc/tap`  
`$ brew install stockuploader`

## Building

Install the most recent version of Go. Currently this is version `1.21.1`. Since all of the dependencies are vendorized, a `go build` is enough. After that you'll see a `stockuploader` binary in the project folder.