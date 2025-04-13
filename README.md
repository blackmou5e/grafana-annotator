# Grafana Annotator

A command-line tool for creating annotations across multiple Grafana dashboards simultaneously.

## Installation

```bash
git clone git@github.com:blackmou5e/grafana-annotator.git
make build
```

## Configuration

Create a configuration file at `$HOME/.gf-annotator/config.yaml` or in the current directory:

```yaml
grafana_host: "localhost"
grafana_port: "3000"
sa_token: "your-token-here"
debug: false
timeout: 30
```

## Usage

Create annotations:
```bash
grafana-annotator create -t "deployment,production" -m "Deployed version 1.2.3"
```

View version information:
```bash
grafana-annotator version
```
