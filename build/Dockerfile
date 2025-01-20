FROM golang:1.23-alpine3.20 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ ./cmd
COPY internal/ ./internal
COPY pkg/ ./pkg

ENV CGO_ENABLED=0
ENV GOOS=linux

RUN go build -o /wallet-backend cmd/app/main.go

FROM alpine:3.20

RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

WORKDIR /app

COPY --from=build /wallet-backend .

RUN mkdir /app/configs
RUN mkdir /app/logs
RUN touch /app/logs/app.log
RUN ln -sf /dev/stdout /app/logs/app.log

EXPOSE 8080

CMD ["./wallet-backend"]
