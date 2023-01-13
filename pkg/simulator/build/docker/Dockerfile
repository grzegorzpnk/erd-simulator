FROM ubuntu:20.04

WORKDIR /nmt

USER root

COPY --chown=root ./bin/simu ./simu
COPY --chown=root ./src/cmd/config.json ./


RUN chmod +x ./simu

ENTRYPOINT ["./simu"]