# Meduza

Command and Control Framework

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
If we are planning to spawn listeners on ports from 8000 to 8010 we can go into the `docker-compose.yml` and add the following line:
```shell
  ...
  teamserver:
    container_name: ${TEAMSERVER_HOSTNAME}
    build:
      context: teamserver
      dockerfile: docker/Dockerfile
    ports:
      - "${TEAMSERVER_PORT}:${TEAMSERVER_PORT}"
      - "${DLV_PORT:-2345}:${DLV_PORT:-2345}"
      - "8000-8010:8000-8010" // Individual ports can be specified as well with "<some_port>:<some_port>"
    ...
```

- To start a listener, a `POST` request should be sent to `http://<server_ip>:<server_port>/api/v1/listeners` with the following body:
```shell
{
        "type":"",
        "name": "",
        "status": ,
        "description": "",
        "config": {}
}
```
which should be modified based on the listener type. 
*Fair notice: the status is an int.*
- After the listener is created, it's UUID can be extracted using a `GET` request to `http://<server_ip>:<server_port>/api/v1/listeners/all`
- The listener can be started using a `POST` request the following endpoint - `http://<server_ip>:<server_port>/api/v1/listeners/<listener_uuid>/start`
- The listener can be stopped using a `POST` request the following endpoint - `http://<server_ip>:<server_port>/api/v1/listeners/<listener_uuid>/stop`
- The listener can be deleted using a `DELETE` request the following endpoint - `http://<server_ip>:<server_port>/api/v1/listeners/<listener_uuid>`
- The listener can be updated using a `PUT` request the following endpoint - `http://<server_ip>:<server_port>/api/v1/listeners/<listener_uuid>`
- The listener can be queried individually using a `GET` request the following endpoint - `http://<server_ip>:<server_port>/api/v1/listeners/<listener_uuid>`

### 6. Starting the Client
Navigate to `Meduza/client` and run the development server:

`npm run dev`
# or
`yarn dev`
# or
`pnpm dev`
# or
`bun dev`


Open [http://localhost:3000](http://localhost:3000) with your browser to see the client.
