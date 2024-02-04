# Simple Load Balancer and Backend Server

This project includes a simple load balancer (`lb`) and a backend server (`backend`) written in Go, which are containerized using Docker. The load balancer forwards incoming HTTP requests to the backend server, which responds with a simple greeting.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

- Go 1.18 or higher
- Docker
- Docker Compose

### Installing

Clone the repository to your local machine:

```bash
git clone https://github.com/yourusername/own-load-balancer.git
```

### Building and Running
Navigate to the root directory of the project and run the following command to build and start the services using Docker Compose:

```bash
docker-compose up --build
```

This command will start two Docker containers:

backend: The backend server listening on port 8080.
lb: The load balancer listening on port 80 and forwarding requests to the backend server.

### Testing

After both services are running, you can send a request to the load balancer using curl:

```bash
curl http://localhost/
```
