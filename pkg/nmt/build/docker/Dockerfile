FROM ubuntu:20.04

WORKDIR /nmt

USER root

COPY --chown=root ./bin/nmt ./nmt
COPY --chown=root ./src/cmd/config.json ./
COPY --chown=root ./src/cmd/mecTopologyforOneCZ.json ./
COPY --chown=root ./src/cmd/networkTopologyforOneCZ.json ./
COPY --chown=root ./src/cmd/mecLinksforOneCZ.json ./

RUN chmod +x ./nmt

ENTRYPOINT ["./nmt"]
