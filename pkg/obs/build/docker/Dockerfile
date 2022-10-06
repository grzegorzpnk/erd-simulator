FROM nicolaka/netshoot

WORKDIR /obs

USER root

COPY --chown=root ./bin/obs ./obs
COPY --chown=root ./src/cmd/config.json ./
COPY --chown=root ./src/cmd/mecTopology.json ./
COPY --chown=root ./src/cmd/networkTopology.json ./

RUN chmod +x ./obs

ENTRYPOINT ["./obs"]