FROM golang:1.19 as base

WORKDIR /app
COPY . /app

RUN go mod download
RUN CGO_ENABLED=0 go build -o main .

FROM gcr.io/distroless/static-debian11:debug

WORKDIR /app

COPY --from=base /app/main .
COPY --from=base /app/logs.txt .

ENTRYPOINT ["/app/main", "logs.txt"]