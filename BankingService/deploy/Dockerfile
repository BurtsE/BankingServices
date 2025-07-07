FROM golang:latest as builder

WORKDIR /app
ADD go.mod .
ADD go.sum .
RUN go mod download

COPY cmd cmd
COPY internal internal
COPY pkg pkg

RUN GOOS=linux go build -o banking_service ./cmd


FROM golang:latest

WORKDIR /app
COPY --from=builder /app/banking_service .
ADD configs/app_config.json configs/app_config.json

EXPOSE 8080

CMD ["./banking_service"]