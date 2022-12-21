FROM ubuntu:20.04

WORKDIR /nmt

USER root

COPY --chown=root ./bin/nmt ./nmt
COPY --chown=root ./src/cmd/config.json ./
COPY --chown=root ./src/cmd/mecTopology.json ./
COPY --chown=root ./src/cmd/networkTopology.json ./
COPY --chown=root ./src/cmd/mecLinks.json ./

RUN chmod +x ./nmt

ENTRYPOINT ["./nmt"]