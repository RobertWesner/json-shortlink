FROM golang:1.26.2-alpine AS build

WORKDIR /build

COPY go.* ./
RUN go mod download

COPY . .
RUN go build -o /out/app .


FROM alpine:latest

WORKDIR /app

COPY --from=build /out/app /app/json-shortlink
COPY .env .env

CMD ["/app/json-shortlink"]
