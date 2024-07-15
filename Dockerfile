FROM golang:1.22-bullseye

WORKDIR /usr/src/app
ADD . /usr/src/app
RUN CGO_ENABLED=0 go build --tags pro -o peeaao-exporter

FROM debian:bullseye-slim
MAINTAINER david@cubiclerebels.com

WORKDIR /usr/src/app
COPY --from=0 /usr/src/app/peeaao-exporter /usr/src/app/peeaao-exporter
RUN apt-get update -y
RUN apt-get install -y curl wget gnupg2
RUN update-ca-certificates

CMD ["/usr/src/app/peeaao-exporter"]
