# CHANGELOG

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/).

## v0.3.0

### Changed

- Using gRPC for communication with gauge

## v0.2.0

### Changed

- Structured logging to support https://github.com/getgauge/gauge/issues/216

## v0.1.1

### Fixed

- Adds a new line character to successful message

## v0.1.0

### Added

- sends KeepAlive pings until report is sent
- Adds zipper module for zipping html-report directory
- Adds sender module for posting html-report archives to `httpserver`
- Adds configurable environment variables