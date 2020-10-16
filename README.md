# Freshping prometheus exporter
[![Build Status](https://github.com/LuxChanLu/freshping-prometheus-exporter/workflows/Build/badge.svg)](https://github.com/LuxChanLu/freshping-prometheus-exporter/actions)
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/55034b0716704ad388252aa3a1789b1a)](https://www.codacy.com/gh/LuxChanLu/freshping-prometheus-exporter/dashboard?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=LuxChanLu/freshping-prometheus-exporter&amp;utm_campaign=Badge_Grade)
[![GitHub release](https://img.shields.io/github/release/LuxChanLu/freshping-prometheus-exporter.svg)](https://github.com/LuxChanLu/freshping-prometheus-exporter/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/LuxChanLu/freshping-prometheus-exporter)](https://goreportcard.com/report/github.com/LuxChanLu/freshping-prometheus-exporter)
[![go-doc](https://godoc.org/github.com/LuxChanLu/freshping-prometheus-exporter?status.svg)](https://pkg.go.dev/github.com/LuxChanLu/freshping-prometheus-exporter)

This is a simple prometheus exporter for [freshping.io](https://www.freshping.io/)

## Configuration

You first need to have a status page enable (This exporter use the status page api)

To run you only need to provide one of these envionment variables :
-  `FRESHPING_URL` : URL to your status page (`https://statuspage.freshping.io/XXXXX-XXXXX/`)
-  `FRESHPING_URL_FILE` : A file containing the URL (In case you want to use a vault with secrets if your url have an url or password `https://username:password@statuspage.freshping.io/XXXXX-XXXXX/`)

## Kubernetes
(You can add the `prometheus.io/scrape` if your deployment use it)

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: freshping-exporter
  labels:
    app.kubernetes.io/name: freshping-exporter
spec:
  selector:
    matchLabels:
      app.kubernetes.io/name: freshping-exporter
  template:
    metadata:
      labels:
        app.kubernetes.io/name: freshping-exporter
    spec:
      containers:
      - name: freshping-exporter
        image: luxchan/freshping-prometheus-exporter:<VERSION>
        resources:
          requests:
            memory: "64Mi"
            cpu: "50m"
          limits:
            memory: "128Mi"
            cpu: "150m"
        ports:
        - containerPort: 9705
          protocol: TCP
          name: metrics
        securityContext:
          allowPrivilegeEscalation: false
          privileged: false
          capabilities:
            drop:
              - ALL
      securityContext:
        readOnlyRootFilesystem: true
        runAsNonRoot: true
        runAsUser: 10001
```

And if needed the `Service` and `ServiceMonitor`
```yaml
---
apiVersion: v1
kind: Service
metadata:
  name: freshping-exporter
  labels:
    app.kubernetes.io/name: freshping-exporter
spec:
  selector:
    app.kubernetes.io/name: freshping-exporter
  ports:
  - port: 9705
    targetPort: 9705
    name: metrics
    protocol: TCP
---
kind: ServiceMonitor
apiVersion: monitoring.coreos.com/v1
metadata:
  name: freshping-exporter
  namespace: monitor
  labels:
    app.kubernetes.io/name: freshping-exporter
spec:
  selector:
    matchLabels:
      app.kubernetes.io/name: freshping-exporter
  jobLabel: app.kubernetes.io/name
  endpoints: 
  - port: "metrics"
    targetPort: 9705
    interval: 15s
```
