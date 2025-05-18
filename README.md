# arbitrage-detector

This service generates dummy forex arbitrage opportunities on a daily basis and exposes an API to fetch the top opportunities for the day. Data is stored in a simple JSON file.

## Building and running

```bash
# Build docker images and run
docker-compose up --build
```

The service will be available on `http://localhost:8080/opportunities/top`.
