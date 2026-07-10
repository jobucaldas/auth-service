from docker.io/golang:1.26.5-alpine3.24 as builder

WORKDIR /app

COPY . .

RUN go mod tidy \
 && go build -o /app/auth .

from docker.io/golang:1.26.5-alpine3.24

WORKDIR /app

RUN addgroup -S auth-service && adduser -S auth-service -G auth-service
USER auth-service

COPY --from=builder --chown=auth-service /app/auth ./auth

CMD ["./auth"]
