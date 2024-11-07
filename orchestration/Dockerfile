FROM golang:1.23 as build

WORKDIR /app

COPY . .

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o orchestration ./cmd

FROM scratch

WORKDIR /app

COPY --from=build /app/orchestration .

COPY cmd/.env .

ENTRYPOINT ["./orchestration"]