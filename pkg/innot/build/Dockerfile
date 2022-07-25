FROM ubuntu:18.04

WORKDIR /innot

USER root

COPY --chown=root ./bin/innot ./
COPY --chown=root ./src/cmd/config.json ./

RUN chmod +x ./innot

ENTRYPOINT ["./innot"]