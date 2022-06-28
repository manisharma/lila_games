FROM golang:alpine as builder
RUN apk update && apk add --no-cache git
WORKDIR /app
COPY . .
RUN go mod download 
WORKDIR /app/cmd
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o game .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/cmd/game .
COPY --from=builder /app/.env .env
EXPOSE 1234
ENTRYPOINT [ "./game" ]