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
go install github.com/nicolekellydesign/webby-api/webby-cli@latest
```

## Usage

Running just `webby-cli` will print the top-level usage, giving a brief overview of all of the commands.

### Setup

All of the database connection info is retrieved from environment variables, specifically:

- WEBBY_DB_USER
- WEBBY_DB_PASSWORD
- WEBBY_DB_NAME
- WEBBY_ROOT

The database schema is created by running `webby-cli init`.

### Users

No admin user is created in the database. Since a valid session token is needed to add a new user, and there are no starting users, a command for adding users is provided: `webby-cli adduser <name> <password>`.
Users created with this command are marked as protected, and cannot be removed via the HTTP API. Protected users can only be removed by using `webby-cli deluser <id>`.
You can view all users with `webby-cli listusers`.

### Serving

Once initial setup is complete, serve the API with `webby-cli serve`.

## License

Copyright &copy; 2021 NicoleKellyDesign

This project is licensed under the terms of the Apache 2.0 license.
