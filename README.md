# Ping Everything Everywhere All At Once (PEEAAO) Exporter

PEEAAO Exporter is a self-hosted Prometheus Exporter that performs a ping on a list of targeted websites and servers using PEEAAO's API Service. It records the response time received from servers around the world as metrics.

For more information on the returned values, please refer to [PEEAAO](https://peeaao.com) documentation.

## Quickstart

To get started, you will need to prepare the following:

1. Get a PEEAAO API Key from [peeaao.com](https://peeaao.com)
2. Have Go 1.22.2 installed to build the binary

Go Build

```
$ go build
```

### Required Environment Variables

| ENV | Optional | Remarks | Example |
| - | - | - | - |
| ADDR | Y | Address:Port that exporter will listen to | localhost:9100
| AUTH_TOKEN | N | AUTH_TOKEN/API_KEY from PEEAAO | |
| LOCATIONS | Y | List of server locations delimited by "," ; Refer to the supported locations at peeaao.com | singapore,nyc |
| TARGETS | N | List of targets delimited by "," | https://peeaao.com |

### Metrics

peeaao-exporter runs on port 9100 and once its running and configured, you can access the metrics at https://localhost:9100/metrics

## Docker Quickstart

```
$ docker build -t peeaao-exporter:latest .
```

```
$ docker run -p 9100:9100 peeaao-exporter:latest
```
