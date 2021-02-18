# Infra
need the two external networks:
```
docker network create --subnet=172.24.255.0/24 mongodb
docker network create --subnet=172.25.255.0/24 nats
docker network create --opt com.docker.network.driver.mtu=1450 dnc
docker network create swpx
```

could add `-d overlay --attachable ` to network command
