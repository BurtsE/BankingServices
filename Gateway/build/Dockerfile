FROM golang:latest as builder

WORKDIR /app
ADD go.mod .
ADD go.sum .
RUN go mod download

COPY cmd cmd
COPY internal internal
COPY pkg pkg
COPY generated generated

RUN GOOS=linux go build -o gateway ./cmd


FROM golang:latest

WORKDIR /app
COPY --from=builder /app/gateway .
ADD configs/app_config.json configs/app_config.json

EXPOSE 8080
EXPOSE 8081

CMD ["./gateway"]