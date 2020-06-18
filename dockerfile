#STEP 1
FROM golang:1.13-alpine AS builder

RUN adduser -D -g 'abaddonuser' abaddonuser

#copy directories
COPY ./src /app

WORKDIR /app


#get the dependencies
RUN apk update && apk add git && apk add ca-certificates && apk add --no-cache make gcc musl-dev linux-headers

#build the binary
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /app/main
RUN chown -R abaddonuser:abaddonuser /app

#STEP 2
FROM alpine:latest

#copy
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder --chown=abaddonuser:abaddonuser /app /app

EXPOSE 8080
USER abaddonuser
VOLUME ["/app/"]
ENTRYPOINT ["/app/main"]