# paperless-uploader
[![Software License](https://img.shields.io/badge/License-MIT-orange.svg?style=flat-square)](https://github.com/gesquive/paperless-uploader/blob/master/LICENSE)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://pkg.go.dev/github.com/gesquive/paperless-uploader)
[![Build Status](https://img.shields.io/circleci/build/github/gesquive/paperless-uploader?style=flat-square)](https://circleci.com/gh/gesquive/paperless-uploader)
[![Coverage Report](https://img.shields.io/codecov/c/gh/gesquive/paperless-uploader?style=flat-square)](https://codecov.io/gh/gesquive/paperless-uploader)
[![Docker Pulls](https://img.shields.io/docker/pulls/gesquive/paperless-uploader?style=flat-square)](https://hub.docker.com/r/gesquive/paperless-uploader)


Watches a directory for files and uploads them to a [paperless-ng](https://github.com/jonaswinkler/paperless-ng) instance.

## Installing

### Compile
This project has only been tested with go1.17+. To compile just run `go install github.com/gesquive/paperless-uploader@latest` and the executable should be built for you automatically in your `$GOPATH`.

Optionally you can clone the repo and run `make install` to build and copy the executable to `/usr/local/bin/` with correct permissions.

### Download
Alternately, you can download the latest release for your platform from [github](https://github.com/gesquive/paperless-uploader/releases).

Once you have an executable, make sure to copy it somewhere on your path like `/usr/local/bin` or `C:/Program Files/`.
If on a \*nix/mac system, make sure to run `chmod +x /path/to/paperless-uploader`.

### Docker
You can also run paperless-uploader from the provided [Docker image](https://hub.docker.com/r/gesquive/paperless-uploader):

```shell
docker run -d -v $PWD/docker:/config -v /path/to/watch:/watch gesquive/paperless-uploader:latest
```

To get the sample config working, you will need to configure the SMTP server and add target configs. 

For more details read the [Docker image documentation](https://hub.docker.com/r/gesquive/paperless-uploader).

### Homebrew
This app is also avalable from this [homebrew tap](https://github.com/gesquive/homebrew-tap). Just install the tap and then the app will be available.
```shell
$ brew tap gesquive/tap
$ brew install paperless-uploader
```

## Configuration

### Precedence Order
The application looks for variables in the following order:
 - command line flag
 - environment variable
 - config file variable
 - default

So any variable specified on the command line would override values set in the environment or config file.

### Config File
The application looks for a configuration file at the following locations in order:
 - `./config.yml`
 - `~/.config/paperless-uploader/config.yml`
 - `/etc/paperless-uploader/config.yml`

Copy `pkg/config.example.yml` to one of these locations and populate the values with your own. Since the config contains a writable API token, make sure to set permissions on the config file appropriately so others cannot read it. A good suggestion is `chmod 600 /path/to/config.yml`.

If you are planning to run this app as a service, it is recommended that you place the config in `/etc/paperless-uploader/config.yml`.

### Environment Variables
Optionally, instead of using a config file you can specify config entries as environment variables. Use the prefix `PAPERLESS_UPLOADER_` in front of the uppercased variable name. For example, the config variable `paperless-url` would be the environment variable `PAPERLESS_UPLOADER_PAPERLESS_URL`.

## Running as a Service
This application was developed to run as a service.

### NIX
You can use upstart, init, runit or any other service manager to run the `paperless-uploader` executable. Example scripts for systemd and upstart can be found in the `pkg/services` directory. A logrotate script can also be found in the `pkg/services` directory. All of the configs assume the user to run as is named `paperless-uploader`, make sure to change this if needed.

### Homebrew
The homebrew tap installs as a service that can be managed with the commands `brew services (start|stop|restart) gesquive/tap/paperless-uploader`. Before running, edit the config file located at `/usr/local/etc/paperless/config.yml`. To debug, read the service logs at `/usr/local/var/log/paperless-uploader.log`.
By default, the service watches the directory `/usr/local/var/paperless-watch`.

## Usage

```console
Watches a directory for files and uploads them to paperless-ng

Usage:
  paperless-uploader [flags]

Flags:
      --config string             Path to a specific config file (default "./config.yml")
  -h, --help                      help for paperless-uploader
  -l, --log-file string           Path to log file (default "/var/log/paperless-uploader.log")
  -t, --paperless-token string    Authenticate the paperless server with this user token
  -u, --paperless-url string      The base URL for your paperless instance
  -x, --upload-path strings       Path to the file(s) to upload, can be entered multiple times or comma delimited.
  -f, --watch-filter string       The inclusive file filter regex for uploads.
  -i, --watch-interval duration   The interval between polling for changes. (default 1s)
  -p, --watch-path string         Directory to watch for files.
      --version                Display the version number and exit
```

Optionally, a hidden debug flag is available in case you need additional output.
```console
Hidden Flags:
  -D, --debug                  Include debug statements in log output
```

## Documentation

This documentation can be found at github.com/gesquive/paperless-uploader

## License

This package is made available under an MIT-style license. See LICENSE.

## Contributing

PRs are always welcome!
