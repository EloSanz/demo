FROM grafana/grafana:latest

USER root

# Aseguramos la existencia de las carpetas internas de Grafana
RUN mkdir -p /var/lib/grafana/dashboards \
             /etc/grafana/provisioning/datasources \
             /etc/grafana/provisioning/dashboards

RUN apk add --no-cache curl sed || true

# 1. Descargamos los dashboards dinámicamente
RUN curl -sL https://grafana.com/api/dashboards/22676/revisions/1/download > /var/lib/grafana/dashboards/fastapi-1.json && \
    curl -sL https://grafana.com/api/dashboards/18739/revisions/1/download > /var/lib/grafana/dashboards/fastapi-2.json

# Parcha el datasource en los JSONs
RUN sed -i 's/${DS_PROMETHEUS}/prometheus/g' /var/lib/grafana/dashboards/*.json

# 2. Generamos datasources.yml dinámicamente dentro del contenedor
RUN printf "apiVersion: 1\ndatasources:\n  - name: Prometheus\n    type: prometheus\n    uid: prometheus\n    url: http://prometheus:9090\n    access: proxy\n    isDefault: true\n  - name: Loki\n    type: loki\n    url: http://loki:3100\n    access: proxy\n" > /etc/grafana/provisioning/datasources/datasources.yml

# 3. Generamos dashboards.yml dinámicamente dentro del contenedor
RUN printf "apiVersion: 1\nproviders:\n  - name: 'FastAPI Dashboards'\n    orgId: 1\n    folder: 'FastAPI'\n    type: file\n    disableDeletion: false\n    updateIntervalSeconds: 10\n    options:\n      path: /var/lib/grafana/dashboards\n" > /etc/grafana/provisioning/dashboards/dashboards.yml

USER grafana
