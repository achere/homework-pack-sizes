# Pack Sizes Calculator

This project provides a web server to calculate optimal pack sizes for orders according to the following rules:
1. Only whole packs can be sent. Packs cannot be broken open.
2. Within the constraints of Rule 1 above, send out the least amount of items to fulfil the order.
3. Within the constraints of Rules 1 & 2 above, send out as few packs as possible to fulfil each
order.

It connects to a user-supplied Postgres database for storing package sizes.

## Getting Started

You can run the application with Docker:

1.  **Clone the repository:**

    ```sh
    git clone https://github.com/achere/homework-pack-sizes.git
    cd homework-pack-sizes
    ```

2.  **Build the image:**

    ```sh
    docker build -t homework-pack-sizes .
    ```
3. **Run the container:**
    ```sh
    docker run -e DB_URL=<postgres_url> -p 8080:8080 homework-pack-sizes
    ```

    The server will start on the port 8080.


## Configuration

You can set the following environment variables using the `-e VAR=<value>` flag of the `docker run` command to configure the application:
- `DB_URL` (required): the connection string for the PostgreSQL database. The format is `postgres://user:password@hostname:port/database_name`.
- `ORDER`: set the default order amount. Set to 251 by default in the Dockerfile.
- `PORT`: set the port for the HTTP server to listen to. Note that you will also need to add port forwarding:
    ```sh
    docker run -e PORT=9090 -p 9090:9090 homework-pack-sizes
    ```


## Running Tests

To run the unit tests for the project (requires Go installed):

```bash
go test ./...
```
