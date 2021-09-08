# Webby-API

[![Go Report Card](https://goreportcard.com/badge/github.com/nicolekellydesign/webby-api)](https://goreportcard.com/report/github.com/nicolekellydesign/webby-api)
[![License](https://img.shields.io/github/license/nicolekellydesign/webby-api.svg)]()

Webby-API is the backend web service for [Webby](https://github.com/nicolekellydesign/webby), our website and CMS. It handles API calls and performs database operations.

## Dependencies

This project uses PostgreSQL as the database of choice. You'll need that set up if you want to use this.

## Building/Installing

Webby-API is built and installed like any other Go project:

```
go build ./cmd/webby-cli
```

Or

```
go install ./cmd/webby-cli
```

## License

Copyright &copy; 2021 NicoleKellyDesign

This project is licensed under the terms of the Apache 2.0 license.
