from alpine:3.14

RUN { \
        echo 'http://mirrors.ustc.edu.cn/alpine/v3.14/main'; \
        echo 'http://mirrors.ustc.edu.cn/alpine/v3.14/community'; \
        echo 'http://mirrors.ustc.edu.cn/alpine/edge/main'; \
        echo 'http://mirrors.ustc.edu.cn/alpine/edge/community'; \
        echo 'http://mirrors.ustc.edu.cn/alpine/edge/testing'; \
    } > /etc/apk/repositories && \
    apk upgrade   && \
    apk add --no-cache tzdata apk-tools rsyslog rsyslog-mysql && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    apk del tzdata
COPY rsyslog.conf /etc/
EXPOSE 514/udp
ENTRYPOINT ["rsyslogd", "-n"]

