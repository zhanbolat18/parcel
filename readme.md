# Parcel Delivery Services

## Repository structure:
* [deliveries](./deliveries) - service to manage deliveries
* [users](./users) - service to manage users and auth logic
* [libs](./libs) - common functions

## Usage

To start a services need `docker`.

`docker compose up -d` - command started all services with PostgreSQL database. After the start, services will be 
available on `localhost`, users on `8080` port, deliveries on `8081` port. Usage ports can see in [docker-compose](./docker-compose.yml) file.

[`http://localhost:8080`](http://localhost:8080) <- users

[`http://localhost:8081`](http://localhost:8081) <- deliveries

Also, after launch, there will be access to Swagger with minimum documentation.

[users swagger](http://localhost:8080/swagger/index.html)

[deliveries swagger](http://localhost:8081/swagger/index.html)