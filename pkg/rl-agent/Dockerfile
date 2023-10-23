# Build stage
FROM python:3.8.12-slim-buster as builder

WORKDIR /app

COPY requirements.txt .

RUN python -m venv venv && \
    . venv/bin/activate && \
    pip install --no-cache-dir -r requirements.txt

COPY get-prediction-schema.json ./json-schemas/
COPY config.json ./config/
COPY rlagent/ .
COPY start-server.sh .
COPY masked-moreVariety2-_20231024-005130.zip model-masked.zip
COPY model-non-masked-1trajectory_ent006_variableload.zip model-basic.zip

RUN chmod +x ./start-server.sh


# Runtime stage
FROM python:3.8.12-slim-buster

WORKDIR /app

COPY --from=builder /app /app

ENV CONFIG_PATH="/app/config/config.json"
ENV GET_PREDICTION_SCHEMA_PATH="/app/json-schemas/get-prediction-schema.json"

EXPOSE 8080

CMD ["/app/start-server.sh"]
