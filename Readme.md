# CatPics API

This Go project allows users to upload, retrieve, update, and delete cat pictures

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

- Go (version 1.22) (if running outside of docker image or for running test suite)
- SQLite3 (if running outside of docker image)
- Any REST client (e.g., Postman, curl) for testing the API endpoints
- Docker

### Installing

1. **Clone the Repository**

    Start by cloning the repository to your local machine.

    ```sh
    git clone https://github.com/sid0jack/catpics-api.git
    ```

2. **Navigate to the Project Directory**

    ```sh
    cd catpics-api
    ```

3. **Build the Docker image**

    This will compile the application and build a docker image.

    ```sh
    sudo docker build --no-cache --progress=plain -t catpics-api .
    ```

4. **Run the Docker image**

    Execute the docker image.

    ```sh
    docker run -dp 8080:8080 -v "$(pwd)/catpics.sqlite3:/root/catpics.sqlite3" catpics-api
    ```

    This will start the server, listening on port 8080.

### Testing the API

You can test the API endpoints using any HTTP client by sending requests to `http://localhost:8080/swagger/index.html#/ followed by the specific endpoint path.

### Running Tests

To run the automated tests for this system, use the following command:

```sh
CGO_ENABLED=1 go test ./...
```
