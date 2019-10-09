FROM golang:1.13 AS builder
WORKDIR /src/app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -mod vendor -o ./build/app
ENTRYPOINT [ "/src/app/build/app" ]

# FROM alpine:3.10
# RUN apk --no-cache add ca-certificates
# COPY --from=builder /src/app/build/app /opt/bin/
# ENTRYPOINT ["/opt/bin/app"]
