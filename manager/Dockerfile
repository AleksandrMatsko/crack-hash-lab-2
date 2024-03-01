FROM golang:1.21
LABEL authors="AMatsko"

WORKDIR /app

COPY ./manager ./manager
COPY ./shared ./shared
RUN cd ./manager && go mod download && go mod verify && CGO_ENABLED=0 make build

ENTRYPOINT ["./manager/manager"]
