#syntax docker/dockerfile:1

FROM golang:1.27.1 as build

WORKDIR /build

RUN mkdir binaries/
COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./

RUN ./scripts/build.sh

FROM debian:13-slim as base
RUN mkdir /home/app
RUN useradd app
RUN chown -R app:app /home/app

USER app
WORKDIR /home/app

COPY --from=build /build/binaries/ ./

EXPOSE 8080

ENTRYPOINT [ "./api" ]
