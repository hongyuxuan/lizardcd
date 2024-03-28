FROM alpine:3.19.1

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories && apk update && apk add sqlite
        
ADD bin/lizardcd-server /opt
WORKDIR /opt
RUN chmod +x /opt/lizardcd-server
EXPOSE 5117
ENTRYPOINT [ "./lizardcd-server" ]