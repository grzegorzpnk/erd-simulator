FROM nicolaka/netshoot

WORKDIR /erc

USER root

COPY --chown=root ./bin/erc ./erc
COPY --chown=root ./config.json ./
COPY --chown=root ./json-schemas ./json-schemas
COPY --chown=root ./ref-schemas ./ref-schemas

RUN chmod +x ./erc

ENTRYPOINT ["./erc"]
