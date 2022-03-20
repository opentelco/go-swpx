FROM golang:1.17 as builder
RUN mkdir -p /src/bin /src/bin/plugins/providers /src/bin/plugins/resources
RUN mkdir -p /swpx/plugins/providers /swpx/plugins/resources
ADD . /src
WORKDIR /src
# build the main app
#RUN go generate
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o /swpx/swpx .

# build provider plugins
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o /swpx/plugins/providers/default ./providers/default
# RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o /swpx/plugins/providers/sait ./providers/sait
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o /swpx/plugins/providers/vx ./providers/vx

# build resource plugins
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o /swpx/plugins/resources/vrp ./resources/vrp
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o /swpx/plugins/resources/raycore ./resources/raycore



FROM alpine:latest
RUN mkdir -p /swpx/config/ /swpx/plugins/providers /swpx/plugins/resources
WORKDIR /swpx/
COPY --from=builder /swpx/ /swpx/
COPY --from=builder /src/config/config-docker.yml /swpx/config/config.yml

CMD ["/swpx/swpx", "start"]
