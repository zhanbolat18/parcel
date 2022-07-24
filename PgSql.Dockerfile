FROM postgres:alpine

COPY multiDb.sh /docker-entrypoint-initdb.d/

RUN chmod +x /docker-entrypoint-initdb.d/multiDb.sh

EXPOSE 5432