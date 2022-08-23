FROM ubuntu:18.04

WORKDIR /obs

USER root

COPY --chown=root ./bin/obs ./
COPY --chown=root ./src/cmd/config.json ./

RUN chmod +x ./obs

ENTRYPOINT ["./obs"]