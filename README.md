# go-web-crontab - Automation runner for scheduled tasks.

This is an automation runner that picks up schedules from
`cron.d/*.cron` and runs scripts referenced in `cron.scripts`. Older
versions of this software were used at [RTV Slovenia](https://rtvslo.si)
for a number of years.

One of the first versions of the front-end just collected the logs of
`crontab -L2` and created a historical view of task execution, including
any performance or stability information (latency, failure rate,
timeouts, ...).

![](docs/assets/webcron-index.png)

![](docs/assets/webcron-detail.png)

## Usage

Currently, no packaging is provided for the project. It's a Go source
code project you can use with `go install`:

```bash
CGO_ENABLED=0 go install github.com/titpetric/go-web-crontab/cmd/webcron@latest
```

The project builds as a static binary which you can include in your own
Docker images. A docker image with the webcron binaries may be provided
in the future, to be included into your *execution environment*.

## About

The project uses a few core packages to bind together functionality:

| Project                                                                  | Usage                                                                 |
|--------------------------------------------------------------------------|-----------------------------------------------------------------------|
| [github.com/go-bridget/mig](https://github.com/go-bridget/mig)           | Database migrations, database model code generation                   |
| [github.com/titpetric/phpscript](https://github.com/titpetric/phpscript) | Front-end phpscript runtime, minimal dashboard layout                 |
| [github.com/titpetric/pdo](https//github.com/titpetric/pdo)              | A request-lived database connection object, type safe generic methods |

The execution log is stored in a sqlite database as defined with
`PLATFORM_DB_CRONTAB` environment variable.

```bash
PLATFORM_DB_CRONTAB=sqlite://webcron.db
```

Written by [Tit Petric](https://titpetric.com).
