FROM alpine:3.19.1

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories && apk update && apk add sqlite

ADD bin/migrate /opt
ADD manifests /opt/manifests
WORKDIR /opt
RUN chmod +x /opt/migrate
ENTRYPOINT [ "./migrate" ]