FROM docker.io/postgres@sha256:78a275d4c891f7b3a33d3f1a78eda9f1d744954d9e20122bfdc97cdda25cddaf as postgres-base

RUN apt update && apt install ca-certificates -y


# https://www.postgresql.org/docs/current/auth-pg-hba-conf.html
COPY pg_hba.conf /config/
COPY certs/postgres.crt /config/postgres.crt
COPY certs/postgres.key /config/postgres.key
COPY certs/root-ca.crt /config/root-ca.crt

RUN chown -R postgres:postgres /config && \
   chmod -R 0600 /config/*  && \
   chown -R postgres:postgres /config && \
   chmod -R 0600 /config/* 