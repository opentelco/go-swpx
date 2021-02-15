FROM golang:latest as builder
RUN mkdir -p /src/bin /src/bin/plugins/providers /src/bin/plugins/resources
ADD . /src
WORKDIR /src
# build the main app
# RUN go generate
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /src/bin/swpx .

# build provider plugins
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /src/bin/plugins/providers/default ./providers/default
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /src/bin/plugins/providers/sait ./providers/sait

# build resource plugins
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /src/bin/plugins/resources/vrp ./resources/vrp
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /src/bin/plugins/resources/raycore ./resources/raycore



FROM alpine:latest  
RUN apk --no-cache add ca-certificates
RUN mkdir -p /swpx/plugins/providers /swpx/plugins/resources
WORKDIR /swpx/
COPY --from=builder /src/bin/swpx /swpx/
COPY --from=builder /src/bin/plugins/providers/default /swpx/plugins/providers/default 
COPY --from=builder /src/bin/plugins/providers/sait /swpx/plugins/providers/sait 
COPY --from=builder /src/bin/plugins/resources/vrp /swpx/plugins/resources/vrp
COPY --from=builder /src/bin/plugins/resources/raycore /swpx/plugins/resources/raycore
COPY --from=builder /src/config/config-docker.yml /swpx/config.yml

CMD ["/swpx/swpx", "start"]
