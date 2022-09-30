FROM nicolaka/netshoot

WORKDIR /obs

USER root

COPY --chown=root ./bin/obs ./obs
COPY --chown=root ./src/cmd/config.json ./

RUN chmod +x ./obs

ENTRYPOINT ["./obs"]