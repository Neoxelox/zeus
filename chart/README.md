# Example Chart

## Installation

`helm upgrade --install --atomic --wait --timeout 60s zeus ./chart/ -f ./chart/values.yaml --namespace olympus --set=image.tag=latest`

## Documentation

> **Access Zeus backend deployed as a ClusterIP service.**

- Run `kubectl port-forward --namespace olympus svc/zeus 1111:1111`
- Go to `localhost:1111`
