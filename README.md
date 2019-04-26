# gauge-reportserver

A simple [Gauge](https://gauge.org/) plugin that will send (POST) the generated `html-report` to a HTTP fileserver such as [gohttpserver](https://github.com/codeskyblue/gohttpserver).

The aim of this plugin is to gather reports from mulitple projects into a `single` place for reference.

[![Gauge Badge](https://gauge.org/Gauge_Badge.svg)](https://gauge.org)

All notable changes to this project are documented in [CHANGELOG.md](CHANGELOG.md).
The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## :hammer_and_pick: Installation

// TODO

## :gear: Configuration

// TODO

### :bulb: Recommendation

// TODO
    - extend plugin kill timeout

## :electric_plug: Running `gohttpserver` locally

Note: Make sure you have `docker` installed.

```bash
docker run -it --rm -p 8000:8000 -v $PWD:/app/public --name gohttpserver codeskyblue/gohttpserver
```

The above should bring up the httpserver on port `8000` at `http://127.0.0.1:8000`

## :wave: Issues & Contributions

Please [open an issue here](../../issues) on GitHub if you have a problem, suggestion, or other comments.

Pull requests are welcome and encouraged! Any contributions should include new or updated tests as necessary to maintain thorough test coverage.

## :scroll: License

This work is licensed under the terms of [GNU Public License version 3.0](http://www.gnu.org/licenses/gpl-3.0.txt)
