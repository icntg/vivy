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

# Redis

```bash
docker run --name redis --rm -d -i -p127.0.0.1:6379:6379 \
-v /home/data/redis:/data redis
```

# rsyslog

```bash
docker run rsyslog/syslog_appliance_alpine help (provides basic info)

docker run --name rsyslog --rm -d -i \
-p127.0.0.1:514:514 \
-v /home/data/rsyslog:/log \
rsyslog/syslog_appliance_alpine

docker run -i -t -d --privileged --name rsyslog --rm \
  -p514:514 -p514:514/udp \
  -v /home/data/rsyslog:/var/log/vivy \
  rsyslog_rsyslog
```