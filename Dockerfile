# syntax=docker/dockerfile:1

# Build the application from source
FROM golang:1.20.5 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

# Download ffmpeg
RUN apt-get update && apt-get -y install tar xz-utils \
    && wget https://johnvansickle.com/ffmpeg/builds/ffmpeg-git-amd64-static.tar.xz \
    && tar xvf ffmpeg-git-amd64-static.tar.xz \
    && mv ffmpeg-git-*-amd64-static ffmpeg \
    && rm ffmpeg-git-amd64-static.tar.xz

RUN make build

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

ENV PATH=$PATH:/ffmpeg

COPY --from=build-stage /app/ffmpeg ffmpeg/
COPY --from=build-stage /app/bin/tgbot-stt-whisper /tgbot-stt-whisper
COPY --from=build-stage /app/.env /.env

USER nonroot:nonroot

ENTRYPOINT ["/tgbot-stt-whisper"]