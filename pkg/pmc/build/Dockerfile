FROM ubuntu:18.04

WORKDIR /pmc

USER root

COPY --chown=root ./bin/pmc ./
COPY --chown=root ./src/cmd/config.json ./

RUN chmod +x ./pmc

ENTRYPOINT ["./pmc"]