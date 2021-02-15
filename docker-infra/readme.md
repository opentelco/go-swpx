# Infra
need the two external networks:
```
docker network create --subnet=172.24.255.0/24 mongodb
docker network create --subnet=172.25.255.0/24 nats
```

could add `-d overlay --attachable ` to network command
