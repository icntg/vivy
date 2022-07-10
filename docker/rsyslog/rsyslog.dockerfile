FROM alpine:3.15

RUN { \
        echo 'http://mirrors4.tuna.tsinghua.edu.cn/alpine/v3.15/main'; \
        echo 'http://mirrors4.tuna.tsinghua.edu.cn/alpine/v3.15/community'; \
        echo 'http://mirrors4.tuna.tsinghua.edu.cn/alpine/edge/main'; \
        echo 'http://mirrors4.tuna.tsinghua.edu.cn/alpine/edge/community'; \
        echo 'http://mirrors4.tuna.tsinghua.edu.cn/alpine/edge/testing'; \
    } > /etc/apk/repositories && \
    apk upgrade   && \
    apk add --no-cache tzdata apk-tools rsyslog rsyslog-mysql logrotate && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    apk del tzdata
COPY *.conf /etc/
EXPOSE 514/udp
ENTRYPOINT ["rsyslogd", "-n"]

