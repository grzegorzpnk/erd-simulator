FROM nicolaka/netshoot

WORKDIR /nmt

USER root

COPY --chown=root ./bin/nmt ./nmt
COPY --chown=root ./src/cmd/config.json ./

RUN chmod +x ./nmt

ENTRYPOINT ["./nmt"]