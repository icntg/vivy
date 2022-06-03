[TOC]

# 依赖

+ MairaDB 数据库
+ Redis 数据缓存与Session
+ (TODO)RSYSLOG 记录日志

# MariaDB

```bash
docker run --name mysql --rm -d -i -p3306:3306 \
-v /home/data/mysql/logs:/logs \
-v /home/data/mysql/data:/var/lib/mysql \
-e MYSQL_ROOT_PASSWORD=root mariadb
```
