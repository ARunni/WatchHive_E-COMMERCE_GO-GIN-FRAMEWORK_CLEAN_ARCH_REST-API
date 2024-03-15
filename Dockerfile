FROM golang:1.21.5-alpine3.18 AS build-stage


WORKDIR /app


COPY ./ /app


RUN mkdir -p /app/build
RUN go mod download
RUN go build -v -o /app/build/api ./cmd/api





EXPOSE 7000


CMD ["/app/build/api"]
