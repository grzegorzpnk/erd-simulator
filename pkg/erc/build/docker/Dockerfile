FROM nicolaka/netshoot

WORKDIR /opt/emco/erc

# RUN addgroup -S emco && adduser -S -G emco emco
# RUN chown emco:emco . -R

COPY --chown=root ./erc ./
COPY --chown=root ./config.json ./
COPY --chown=root ./json-schemas ./json-schemas
COPY --chown=root ./ref-schemas ./ref-schemas

USER root

ENTRYPOINT ["./erc"]
