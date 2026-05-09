# Docker Commands

```bash
docker pull mysql:8.4

docker run --name ecomm-mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=password -d mysql:8.4
```

```bash
docker ps

docker stop <container_id>

docker ps -a
docker rm <container_id>

docker rmi <image_id>
```
