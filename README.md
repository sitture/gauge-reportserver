# gauge-report-server

A simple gauge plugin that will send/post the generated html-report to a HTTP fileserver such as [gohttpserver](https://github.com/codeskyblue/gohttpserver).

The aim of this plugin is to gather reports from mulitple projects into a `single` place for reference.

## Configuration

// TODO

## Running `gohttpserver` locally

Note: Make sure you have `docker` installed.

```bash
docker run -it --rm -p 8000:8000 -v $PWD:/app/public --name gohttpserver codeskyblue/gohttpserver
```

The above should bring up the httpserver on port `8000` at `http://127.0.0.1:8000`

## License

// TODO