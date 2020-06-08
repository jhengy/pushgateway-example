# Pushgateway example

This is a code example for the blog post 
["Monitoring Batch Jobs with Prometheus"](https://www.nicktriller.com/blog/monitoring-batch-jobs-with-prometheus/).

## Usage

```
docker-compose up -d
```

The Prometheus web UI should be accessible on [http://localhost:9090](http://localhost:9090).

# Aim
- illustrate how to export an metrics from an instrumented golang program to pushgateway which is to be configured to be scraped by prometheus
    - how to export metrics in instrumented golang program to pushgateway
    - how to configure pushgateway to receive and send metrics
    - how to configure prometheus receiver to receive metrics