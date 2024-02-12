FROM golang:1.20-alpine3.16 as backend-build

WORKDIR /app

RUN apk update && apk add --no-cache gcc musl-dev

COPY . .

RUN go build -o main .

FROM alpine:3.16

WORKDIR /app

COPY --from=backend-build /app/main .

COPY dev_config.json .

COPY /database/database.db ./database/database.db

COPY key.pem cert.pem ./

EXPOSE 8080

CMD ["./main"]