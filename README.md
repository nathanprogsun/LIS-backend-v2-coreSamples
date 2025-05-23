# LIS-backend-v2-coreSamples

## Intro
For details about coresamples v2 service, see documents under [LIS-Coresamples-v2](https://vibrantamerica.atlassian.net/wiki/spaces/LIS/pages/707559427/LIS-Coresamples-v2)

---

## Docker Compose Local Development Setup

This document describes how to run the LIS CoreSamples service and its core dependencies using Docker Compose for local development and testing.

### Prerequisites

*   **Docker:** Ensure Docker Engine is installed and running. (Refer to official Docker documentation for installation)
*   **Docker Compose:** Ensure Docker Compose (usually included with Docker Desktop, or as a standalone plugin) is installed. (Refer to official Docker Compose documentation)
*   **Git:** To clone the repository.
*   **Make (Optional but Recommended):** The project uses a Makefile for common tasks.
*   **Go (Optional):** If you need to modify Go code and rebuild, Go 1.19+ is required locally. The Docker build process handles Go compilation within a container.

### Project Structure (Relevant for Docker Compose)

*   `Dockerfile`: Used to build the `coresamples` service image.
*   `docker-compose.yaml`: Defines the services, networks, and volumes for the local environment.
*   `.env` (Create this file): Used to store sensitive credentials and local configuration overrides. **Do not commit this file.**
*   `.env.example` (Recommended to create): An example file showing the variables needed in `.env`.
*   `Makefile`: Contains helper targets like `make tidy`, `make build`.

### Setup and Running

1.  **Clone the Repository:**
    ```bash
    # git clone <repository_url>
    # cd <repository_directory>
    ```
    (Assuming you are already in the cloned directory)

2.  **Configuration (`.env` file):**

    Create a file named `.env` in the root of the project directory (where `docker-compose.yaml` is located). This file will store your local configurations and secrets. **Do not commit the `.env` file to version control.**

    **Example `.env` file content (save as `.env`):**
    ```env
    # MySQL Credentials
    MYSQL_ROOT_PASSWORD_VAL=myrootpassword_changeme
    MYSQL_USER_VAL=coresamples_user
    MYSQL_PASSWORD_VAL=coresamples_pass_changeme
    MYSQL_DATABASE_VAL=coresamples_db

    # JWT Secret for coresamples service
    JWT_SECRET_VAL=thisisadevelopmentsecret_pleasedontuseinprod_changeme

    # Optional: Override other defaults from docker-compose.yaml if needed
    # CORESAMPLES_LOG_LEVEL=info
    # SENTRY_DSN_VAL=your_sentry_dsn_here

    # Optional: Override default host port mappings if they conflict
    # MYSQL_HOST_PORT=3308
    # CONSUL_HOST_PORT=8501
    # REDIS_HOST_PORT=6380
    # JAEGER_AGENT_UDP_PORT=6832 # Host port for Jaeger agent
    # JAEGER_UI_HOST_PORT=16687
    ```
    The `docker-compose.yaml` is configured to require some of these variables (e.g., `JWT_SECRET_VAL`, database credentials). If they are not set, Docker Compose will show an error.

3.  **Build and Start Services:**
    Open a terminal in the project root directory and run:
    ```bash
    docker-compose up --build
    ```
    *   `--build`: Forces Docker Compose to build the `coresamples` image using the `Dockerfile`. Required on the first run or after code changes.
    *   Remove `--build` for subsequent runs if you haven't changed the `coresamples` source code or `Dockerfile`.
    *   To run in detached mode (in the background), add the `-d` flag: `docker-compose up --build -d`.

4.  **Accessing Services:**

    *   **CoreSamples HTTP API:** `http://localhost:8083`
    *   **CoreSamples gRPC API:** `localhost:8084` (requires a gRPC client)
    *   **Consul UI:** `http://localhost:${CONSUL_HOST_PORT:-8500}` (default: `http://localhost:8500`)
    *   **Jaeger UI:** `http://localhost:${JAEGER_UI_HOST_PORT:-16686}` (default: `http://localhost:16686`)
    *   **MySQL:** Accessible on `localhost:${MYSQL_HOST_PORT:-3307}` (default: `localhost:3307`) via a MySQL client.
        *   User: Value of `MYSQL_USER_VAL` from your `.env` file.
        *   Password: Value of `MYSQL_PASSWORD_VAL` from your `.env` file.
        *   Database: Value of `MYSQL_DATABASE_VAL` from your `.env` file.
    *   **Redis:** Accessible on `localhost:${REDIS_HOST_PORT:-6379}` (default: `localhost:6379`) via `redis-cli`.

5.  **Viewing Logs:**
    If running in detached mode, or to view logs from all services:
    ```bash
    docker-compose logs -f
    ```
    To view logs for a specific service (e.g., `coresamples`):
    ```bash
    docker-compose logs -f coresamples
    ```

6.  **Stopping Services:**
    To stop the services, press `Ctrl+C` in the terminal where `docker-compose up` is running.
    If running in detached mode, or from another terminal:
    ```bash
    docker-compose down
    ```
    To stop and remove volumes (warning: this deletes data like your database):
    ```bash
    docker-compose down -v
    ```

### Development Workflow

1.  **Modify Go code** in the `coresamples` project.
2.  **Rebuild and restart** the `coresamples` service:
    ```bash
    docker-compose up --build -d coresamples
    ```
    Or, if you prefer to restart all services:
    ```bash
    docker-compose down && docker-compose up --build -d
    ```

### Important Notes & Limitations

*   **Go Application Adaptation:** This Docker Compose setup assumes the `coresamples` Go application has been adapted (or will be adapted) to:
    *   Recognize the `RUN_ENV=dev_docker_compose` environment variable.
    *   Prioritize other environment variables (e.g., for DB, Redis, JWT_SECRET) over Consul lookups for critical connection settings.
    *   Connect to Redis as a standalone instance.
    *   Handle potential Kafka unavailability gracefully (Kafka is not included in this setup).
*   **Kafka:** Apache Kafka is **not** included in this Docker Compose setup to keep it lightweight.
*   **External Microservices:** Other external microservices that `coresamples` might interact with are not part of this setup.
*   **Resource Usage:** Running multiple services can be resource-intensive.
*   **Security:** The `.env` file is for local development convenience. Ensure it's not committed and that production configurations use secure secret management.

### Troubleshooting

*   **Missing Environment Variables:** If Docker Compose fails with an error like `Variable is not set and no default was provided`, ensure the required variable (e.g., `JWT_SECRET_VAL`) is defined in your `.env` file.
*   **Service Connection Issues:** Verify hostnames and ports in `coresamples` environment variables match Docker Compose service names.
*   **Port Conflicts:** Use the `*_HOST_PORT` variables in your `.env` file (e.g., `MYSQL_HOST_PORT=3308`) to change host-side port mappings if defaults conflict with local services.
*   **Build Failures:** Check logs from `docker-compose build`.
*   **Data Persistence:** Data is stored in Docker named volumes. Use `docker-compose down -v` for a clean start.