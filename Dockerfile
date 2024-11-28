FROM golang:1.23 as builder
ENV GOPROXY=https://goproxy.cn,direct
WORKDIR /app
COPY .netrc /root/
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY classin/etc/classin-api.yaml /app/etc/config.yaml
RUN go build -ldflags="-s -w" -o /app/lesson classin/classin.go
FROM alpine:latest
ENV ZONEINFO=/usr/share/zoneinfo
RUN apk --no-cache add libc6-compat libgcc libstdc++
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /app/lesson /usr/local/bin/lesson
COPY --from=builder /app/etc/config.yaml /etc/config.yaml
EXPOSE 8888
CMD ["lesson", "-f", "/etc/config.yaml"]