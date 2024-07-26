FROM golang:1.22.5-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /project-api ./main.go

FROM alpine:latest

COPY --from=builder /project-api /project-api

EXPOSE 8080

COPY .env /app/.env

CMD ["/project-api"]
