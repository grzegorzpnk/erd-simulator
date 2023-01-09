FROM ubuntu:20.04

WORKDIR /erc

USER root

COPY --chown=root ./bin/erc ./erc
COPY --chown=root ./src/cmd/config.json ./
COPY --chown=root ./json-schemas ./json-schemas
COPY --chown=root ./ref-schemas ./ref-schemas

RUN chmod +x ./erc

ENTRYPOINT ["./erc"]
