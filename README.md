# Eve Data Service

Proof of concept application responsible for managing EVE Online Data.

[EVE Online API](https://esi.evetech.net/ui/#/)

## Prerequisites

- Go 1.22
- Taskfile [https://taskfile.dev/](https://taskfile.dev/)
- Docker

## Testing

- Ensure that docker is running.

`task test`

A file called `coverage.html` will be generated that shows the test coverage of the code.

## Building

- Ensure that docker is running.

`task build`

## Running

- Ensure that docker is running.

`task run`

## Tasks

To list all available tasks run:
`task -l`

## Configuration

The application is configured using [dotenvconfig](https://github.com/andrewapj/dotenvconfig)

- By default, the file `local.env` is used to configure the application.
- An environment variable `APP_CONFIG_PATH` can be set to override this.

Configuration files should end with `.env` and be in the following format:
`TEST_KEY=123`