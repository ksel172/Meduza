# Meduza

Command and Control Server.

## Development

### 1. Configuration

Run the following script to create `.env.development` and `.env.production`:

```bash
grep -E '^\s+[A-Z_]+:' docker-compose.yml | sed -E 's/^\s+([A-Z_]+):\s*(.*)/\1=\2/' | tee .env.development .env.production
```

Use `.env.example` to see an example configuration.

### 2. Running Services

To run the services, we use Docker compose:

```shell
docker compose --env-file .env.development up --force-recreate --build
```

To delete volumes in case the database needs to be reinitialized:

```shell
docker compose --env-file .env.development down --volumes
```

This will build the application server, create an instance of a PostgresSQL database and launch Postgres Admin web app.
