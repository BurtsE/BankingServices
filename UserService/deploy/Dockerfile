FROM golang:latest as builder

WORKDIR /app
ADD go.mod .
ADD go.sum .
RUN go mod download

COPY cmd cmd
COPY internal internal

RUN GOOS=linux go build -o user_service ./cmd


FROM golang:latest

WORKDIR /app
COPY --from=builder /app/user_service .
ADD configs/app_config.json configs/app_config.json

EXPOSE 8080

CMD ["./user_service"]