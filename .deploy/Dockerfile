FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod ./
RUN go mod download && go mod tidy
COPY . .
RUN go mod tidy
RUN go build -o isaac-fansite .

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/isaac-fansite .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static

CMD ["./isaac-fansite"]