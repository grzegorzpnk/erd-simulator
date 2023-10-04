FROM ubuntu:20.04

WORKDIR /nmt

USER root

COPY --chown=root ./bin/nmt ./nmt
COPY --chown=root ./src/cmd/config.json ./
COPY --chown=root ./src/cmd/mecTopologyforTwoCZ.json ./
COPY --chown=root ./src/cmd/networkTopologyforTwoCZ.json ./
COPY --chown=root ./src/cmd/mecLinksforTwoCZ.json ./

RUN chmod +x ./nmt

ENTRYPOINT ["./nmt"]
