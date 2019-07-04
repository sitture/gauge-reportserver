# gauge-reportserver

A simple [Gauge](https://gauge.org/) plugin that will send (POST) the generated `html-report` to a HTTP fileserver such as [gohttpserver](https://github.com/codeskyblue/gohttpserver).

The aim of this plugin is to gather reports from mulitple projects into a `single` place for reference.

[![Gauge Badge](https://gauge.org/Gauge_Badge.svg)](https://gauge.org)

All notable changes to this project are documented in [CHANGELOG.md](CHANGELOG.md).
The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## :hammer_and_pick: Installation

* Install the plugin
* `gohttpserver` running somewhere, or refer to [Running gohttpserver locally](#electric_plug-running-gohttpserver-locally)

```sh
gauge install reportserver
```

### Offline installation

* Download the plugin from [Releases](../../releases)

```sh
gauge install reportserver --file reportserver-${version}-darwin.x86_64.zip
```

### Using the plugin

Add `reportserver` to your project's `manifest.json`.

E.g.

```json
{
  "Language": "java",
  "Plugins": [
    "html-report",
    "reportserver"
  ]
}
```

## :gear: Configuration

You can set the following environment variables to override the configuration OR by adding these to `env/default.properties`:

- `REPORTSERVER_HOST` - This is the base url of the http server. Default is set to `http://localhost:8000`
- `REPORTSERVER_BASE_DIR` - This is the base directory of your reports. Default is set to your project directory name.
- `REPORTSERVER_PATH` - This is path where you want the report files to go. if this is not specified, then the environment directory name is used as the path.
- `REPORTSERVER_TIMEOUT_IN_SECONDS` - This is how long to wait for html-report to be ready before sending. Default is 15 seconds.

Examples:

```sh
REPORTSERVER_HOST=http://myreportserver.com
REPORTSERVER_BASE_DIR=myproject

# Path on reportserver
http://myreportserver.com/myproject/${env_directory}

REPORTSERVER_HOST=http://myreportserver.com
REPORTSERVER_BASE_DIR=myproject
REPORTSERVER_PATH=test/test

# Path on reportserver
http://myreportserver.com/myproject/test/test/
```

## :electric_plug: Running `gohttpserver` locally

Note: Make sure you have `docker` installed.

```bash
docker run -it --rm -p 8000:8000 -v $PWD:/app/public --name gohttpserver codeskyblue/gohttpserver
```

* You can also use `docker-compose` to bring up the service. Create a new file `docker-compose.yml` and add the following:

```sh
version: '2'
services:
  gohttpserver:
    image: codeskyblue/gohttpserver
    ports:
      - '8000:8000'
    volumes:
      - '.:/app/public'
```

Run `docker-compose up -d` to bring up the gohttpserver in background.

The above should bring up the httpserver on port `8000` at `http://127.0.0.1:8000`

## Building locally

```bash
go run build/make.go
go run build/make.go --install
```

## :wave: Issues & Contributions

Please [open an issue here](../../issues) on GitHub if you have a problem, suggestion, or other comments.

Pull requests are welcome and encouraged! Any contributions should include new or updated tests as necessary to maintain thorough test coverage.

## :scroll: License

This work is licensed under the terms of [GNU Public License version 3.0](http://www.gnu.org/licenses/gpl-3.0.txt)
