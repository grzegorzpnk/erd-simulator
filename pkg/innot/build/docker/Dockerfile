FROM nicolaka/netshoot

WORKDIR /innot

USER root

COPY --chown=root ./bin/innot ./innot
COPY --chown=root ./src/cmd/config.json ./

RUN chmod +x ./innot

ENTRYPOINT ["./innot"]