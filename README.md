# Meduza

Command and Control Server.

## Development

### 1. Configuration

See [`.env.example`](.env.example) for an example configuration.

### 2. Running Services

To run the services, we use Docker compose:

```shell
docker compose --env-file .env.dev up --force-recreate --build
```

This will build the application server, create an instance of a PostgresSQL and Redis databases and launch Postgres Admin web app.

#### Run Mode

The `TEAMSERVER_MODE` environmental variable can be used to control whether to run the server with a Delve debugger or without.

In case you're using the `TEAMSERVER_MODE=debug`, configure `DLV_PORT` env var and set up the Delve debugger client.
[Available Delve clients](https://github.com/go-delve/delve/blob/master/Documentation/EditorIntegration.md).

### 3. Cleaning Up
To delete volumes in case the database needs to be recreated:

```shell
docker compose --env-file .env.development down --volumes
```

In some cases, you will also need to delete the database from the filesystem using:

```shell
docker volumes rm $VOLUME_NAME
```


