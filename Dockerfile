FROM alpine:latest

COPY ./ ./

# install psql
RUN apk add --no-cache --update \
    postgresql-client \
    ca-certificates

RUN chmod -R +x ./scripts

CMD ["./scripts/prepare-database.sh ./migrations ./scripts/_init.sql db ./server"]
