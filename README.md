# Webby-API

[![Go Report Card](https://goreportcard.com/badge/github.com/nicolekellydesign/webby-api)](https://goreportcard.com/report/github.com/nicolekellydesign/webby-api)
[![License](https://img.shields.io/github/license/nicolekellydesign/webby-api.svg)]()

Webby-API is the backend web service for [Webby](https://github.com/nicolekellydesign/webby), our website and CMS. It handles API calls and performs database operations.

## Dependencies

This project uses PostgreSQL as the database of choice. You'll need that set up if you want to use this.

## Building

To generate authentication tokens, Webby-API uses a variable that is injected during the build process. The easiest way (that I know of so far) to do this is to have a file in the root of the project named `.signingKey` (this name is in the `.gitignore` already) with the text of the key to use.

To generate a key, you can run `openssl rand -base64 172`

Once you have a key, you can inject it using Go's ldflags like so:
```
export SIGNING_KEY=$(cat .signingKey) && go build -ldflags "-X github.com/nicolekellydesign/webby-api/server/signingKey=$SIGNING_KEY" ./cmd/webby-cli
```

If you have another way to inject the variable that you like better, then feel free to do that.

## License

Copyright &copy; 2021 NicoleKellyDesign

This project is licensed under the terms of the Apache 2.0 license.
