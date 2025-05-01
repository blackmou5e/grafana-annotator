# Grafana Annotator

A command-line tool for creating annotations across multiple Grafana dashboards simultaneously.

## Installation

```bash
git clone git@github.com:blackmou5e/grafana-annotator.git
make build
```

## Configuration

Create a configuration file at one of this paths:
* `$HOME/.gf-annotator/config.yaml`
* `$XDG_CONFIG_HOME/grafana-annotator/config.yaml`
* `./grafana_annotator.yaml`


```yaml
annotator_grafana_url: "http://localhost:3000"
annotator_sa_token: "your-token-here"
annotator_debug: false
annotator_timeout: 30
```

Or use environment variables:
```bash
export ANNOTATOR_GRAFANA_URL='http://localhost:3000'
export ANNOTATOR_SA_TOKEN='your-token-here'
export ANNOTATOR_DEBUG=false
export ANNOTATOR_TIMEOUT=30
```

In both examples above, grafana url, debug and timeout are default values, that will be used if no config provided.

## Usage

Create annotations:
```bash
grafana-annotator create -t "deployment,production" -m "Deployed version 1.2.3"
```

View version information:
```bash
grafana-annotator version
```
