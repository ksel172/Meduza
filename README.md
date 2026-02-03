# Meduza

Meduza is a modular, collaborative C2 framework written in Go and Docker. The [Meduza Framework](https://github.com/Meduza-Framework) also features a client written in React and an agent written in C#.

## Features
W.I.P

## Quick Start

A detailed guide on installation, configuration and usage as well as the project's architecture and development information can be found in the official [Meduza Documentation](https://meduza-framework.github.io/meduza-documentation/)

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

### 4. JWT Secret
In case you want generate a JWT secret, run the following command:
```bash
openssl rand -base64 64
```
If openssl is not install on your device, install it first based on operating system first.
After generating JWT secret, add it your .env file as follows:
```bash
JWT_SECRET=your_generated_secret
```


### 5. Listener creation

WIP:
Before creating a listener, an array of ports that will be opened in the docker container of the teamserver should be specified. 
That can be done in the `.env` file using the `LISTENER_PORT_RANGE_START` and `LISTENER_PORT_RANGE_END` variables.

In order to prevent unauthorized users from exploiting the API endpoints, it's crucial to implement authentication and authorization checks. Ensure that all endpoints specified for listener creation, manipulation, and querying are secured and only accessible to authenticated users. Use JWT tokens for this purpose, validating them on each request. Update the documentation to specify that JWT authentication is required and provide an example token acquisition method.
### 6. Starting the Client
Navigate to `Meduza/client` and run the development server:

```bash
npm run dev
# or
yarn dev
# or
pnpm dev
# or
bun dev
```

Open [http://localhost:3000](http://localhost:3000) with your browser to see the client.
