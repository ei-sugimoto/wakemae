# Wakemae

Wakemae is a lightweight DNS server for Docker containers. Simply add labels to your containers to automatically access them by domain name.
This project is similar dnsdock.

## Features

- **Automatic DNS Registration**: Monitors Docker container start/stop events and automatically manages DNS records
- **Label-based Configuration**: Uses `wakemae.domain` labels to set domain names
- **Real-time Updates**: Monitors and reflects container state changes in real-time
- **Fallback Support**: Forwards unregistered domains to upstream DNS servers
- **A/CNAME Record Support**: Supports both A records and CNAME records



## Usage
### Docker Compose Example (reccomend)

Sample configuration is available in `example/docker-compose.yml`:

```yaml
services:
  wakemae:
    image: ghcr.io/ei-sugimoto/wakemae:latest
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
    networks:
      wakemae_network:
        ipv4_address: 172.20.0.10
  
  web1:
    image: nginx:latest
    ports:
      - 8080:80
    labels:
      wakemae.domain: web1.docker
    dns:
      - 172.20.0.10
    networks:
      wakemae_network:
        ipv4_address: 172.20.0.11
    depends_on:
      - wakemae

networks:
  wakemae_network:
    driver: bridge
    ipam:
      config:
        - subnet: 172.20.0.0/16
```

With this configuration, you can access the nginx container via `web1.docker`.

```bash
docker compose exec web1 bash -c "apt install curl && curl http://web2.docker/"
```

## Architecture

Wakemae consists of the following components:

- **DNS Server**: Handles DNS queries on both UDP and TCP
- **Docker Listener**: Monitors Docker events to detect container start/stop
- **Registry**: Manages domain name to IP address mappings

### Flow

1. When a Docker container starts, the Docker Listener detects the event
2. Checks the container's `wakemae.domain` label
3. If the label exists, registers the IP address and domain name mapping in the Registry
4. When DNS queries arrive, references the Registry to return IP addresses
5. When containers stop, removes the corresponding records

## Commands

### serve

Start the DNS server and Docker monitoring.

```bash
wakemae serve
```

### Options

In the current version, the following settings are hardcoded:

- DNS Port: 53 (UDP/TCP)
- Bind Address: 0.0.0.0:53
- Upstream DNS Server: 8.8.8.8:53

## Development
### Running Tests

```bash
make test
```

### Running Linter

```bash
make lint
```

### Building

```bash
go build -o wakemae .
```

## License

This project is released under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contributing

Pull requests and issue reports are welcome. Please check existing issues before contributing.

## Notes
- Wakemae requires access to the Docker socket
- Thoroughly test before using in production environments
- DNS port 53 may require root privileges 