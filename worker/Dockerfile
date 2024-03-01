FROM golang:1.21
LABEL authors="AMatsko"

WORKDIR /app

COPY ./worker ./worker
COPY ./shared ./shared
RUN cd ./worker && go mod download && go mod verify && CGO_ENABLED=0 make build

ENTRYPOINT ["./worker/worker"]
