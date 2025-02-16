FROM golang:1.23.6-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/merchstore ./cmd/app/main.go

FROM alpine:latest
RUN apk --no-cache add docker-cli
WORKDIR /app

COPY --from=build /app/merchstore /app/merchstore
COPY --from=build /app/config/application.yml /app/

RUN mkdir -p /app/logs

CMD ["/app/merchstore"]