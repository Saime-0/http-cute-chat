FROM alpine:latest

COPY ./ ./
RUN chmod -R +x ./scripts

# install packages
RUN apk add --no-cache --update \
    postgresql-client \
    ca-certificates \
    git

# install golang
COPY --from=golang:1.17.3-alpine3.15 /usr/local/go/ /usr/local/go/
ENV GOPATH=$HOME/go
ENV PATH=/usr/local/go/bin:$GOPATH/bin:$PATH


# install migrator
RUN ./scripts/install-migrator.sh
#ENV PATH="/usr/local/migrator:${PATH}"
ENV PATH=./.tmp/migrator:$PATH

# build server
RUN go build -v ./server.go

CMD ["./scripts/prepare-database.sh ./migrations ./scripts/_init.sql db ./server"]
