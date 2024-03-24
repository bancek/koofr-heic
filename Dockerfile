FROM golang:1.22.1 AS build

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go install github.com/revel/cmd/revel
RUN revel build . /koofr-heic

FROM ubuntu:22.04

RUN apt-get update && \
  apt-get install -y ffmpeg imagemagick ca-certificates && \
  update-ca-certificates && \
  apt-get autoremove -y && \
  rm -rf /var/lib/apt/lists/*

COPY --from=build /koofr-heic /koofr-heic

CMD /koofr-heic/run.sh
