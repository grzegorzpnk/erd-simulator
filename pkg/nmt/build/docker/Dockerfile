FROM nicolaka/netshoot

WORKDIR /nmt

USER root

COPY --chown=root ./bin/nmt ./nmt
COPY --chown=root ./src/cmd/config.json ./
COPY --chown=root ./src/cmd/mecTopology.json ./
COPY --chown=root ./src/cmd/networkTopology.json ./

RUN chmod +x ./nmt

ENTRYPOINT ["./nmt"]