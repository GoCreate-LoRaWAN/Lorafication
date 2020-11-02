# lorafication daemon

This daemon is a notification service that can be utilized within the infrastructure of LoRa-WAN.

## Table of Contents

- [Running](#running)
    - [Dependencies](#dependencies)
    - [Configuration](#environment-variables)
        - [From Environment](#from-environment)
        - [From File](#from-environment)
    - [Make Rules](#make-rules)

## Running

### Dependencies

The only dependencies to run the services in this repository are:

- `docker`
- `docker-compose`

### Configuration

The lorafication daemon can be configured either via environment (default), or by file. If you want to configure the
daemon via a file, you must pass the `-envFile=` flag when running the daemon, giving it the path to either a JSON or
YAML file containing configuration values.

#### From Environment

The program looks for the following environment variables:

- `LORAFICATION_PORT`: The port that the lorafication daemon listens to/serves from (Default: `9000`).
- `LORAFICATION_LOG_LEVEL`: The log level for the lorafication daemon to use when logging, effectively filtering some
statements. See the [zap documentation](https://godoc.org/go.uber.org/zap/zapcore#Level) here for allowed levels
(Default: `0` (`InfoLevel`)).
- `LORAFICATION_DB_USER`: The postgres database username that gets used within the postgres connection string (Default:
`root`).
- `LORAFICATION_DB_PASS`: The postgres database password that gets used within the postgres connection string (Default:
`root`).
- `LORAFICATION_DB_NAME`: The postgres database name that gets used within the postgres connection string (Default:
`lorafication`).
- `LORAFICATION_DB_HOST`: The postgres database host name that gets used within the postgres connection string (Default:
`db`).
- `LORAFICATION_DB_PORT`: The postgres database port that gets used within the postgres connection string (Default:
`5432`).
- `LORAFICATION_SMTP_HOST`: The SMTP server's host address to connect to in order to send emails
(Default: `smtp.gmail.com`).
- `LORAFICATION_SMTP_PORT`: The SMTP server's port to use in conjunction with the host address to connect to in order to
send emails (Default: `587`)
- `LORAFICATION_SMTP_USER`: The username to use when connecting to the SMTP server (Default: n/a).
- `LORAFICATION_SMTP_PASS`: The password to use when connecting to the SMTP server (Default: n/a).
- `LORAFICATION_READ_TIMEOUT`: The time of the read timeout of any outgoing read requests made by the internal HTTP
server (Default: `10s`).
- `LORAFICATION_WRITE_TIMEOUT`: The time of the read timeout of any outgoing write requests made by the internal HTTP
server (Default: `20s`).
- `LORAFICATION_SHUTDOWN_TIMEOUT`: The time of the graceful shutdown timeout of the lorafication daemon. This is the
amount of time in between an attempted, non-forceful shutdown and the finishing of open requests and/or the shutdown of
integrated services, such as the database (Default: `20s`).

If the environment variable has a supplied default and none are set within the context of the host machine, then the
default will be used.

#### From File

If you opt to set the daemon's configuration via a file, the following structures will need to be used for JSON or YAML,
respectively (in the following examples all fields will be set to their defaults, for granular descriptions of these
fields, see their equivalents in [From Environment](#from-environment)):

JSON:
```json
{
    "port": 9000,
    "logLevel": 0,
    "dbUser": "root",
    "dbPass": "root",
    "dbName": "lorafication",
    "dbHost": "db",
    "dbPort": 5432,
    "smtpHost": "smtp.gmail.com",
    "smtpPort": 587,
    "smtpUser": "<no default>",
    "smtpPass": "<no default>",
    "readTimeout": "10s",
    "writeTimeout": "20s",
    "shutdownTimeout": "20s"
}
```

YAML:
```yaml
port: 9000
logLevel: 0
dbUser: root
dbPass: root
dbName: lorafication
dbHost: db
dbPort: 5432
smtpHost: smtp.gmail.com
smtpPort: 587
smtpUser: <no default>
smtpPass: <no default>
readTimeout: 10s
writeTimeout: 20s
shutdownTimeout: 20s
```

### Make Rules

To run the services simply execute the following command:

```shell
make run
```

This will stop any containers defined by the compose file if already running and then rebuild the containers using the
compose file. The lorafication daemon (`loraficationd`) will be available at `localhost:9000` and the postgres instance
will be available at `localhost:5432`.

To retrieve the logs of the services:

```shell
make logs
```

To stop the services:

```shell
make stop
```

To remove all networks and volumes that the services have created (effectively allowing you to start from a clean state
upon the next run):

```shell
make down
```
