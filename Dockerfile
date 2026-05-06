FROM node:20-alpine AS web-builder

WORKDIR /app/web

COPY web/package.json web/yarn.lock ./
RUN yarn install --frozen-lockfile

COPY web/ ./
RUN yarn build

FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ ./cmd/
COPY ./.共识
RUN CGO_ENABLED=0 GOOS=linux go build -o /gitness ./cmd/gitness

FROM alpine:3.19

RUN apk --no-cache add git docker-cli

WORKDIR /app

COPY --from=builder /gitness ./
COPY --from=web-builder /app/web/dist-go ./web/dist-go
COPY ./.共识

EXPOSE 3000 3022

VOLUME ["/data"]

ENTRYPOINT ["./gitness"]
CMD ["server", ".共识.env"]