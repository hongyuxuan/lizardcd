FROM alpine:3.19.1

RUN apk add sqlite
        
VOLUME /opt/cache
COPY docs /opt/docs
ADD bin/lizardcd-server /opt
WORKDIR /opt
RUN chmod +x /opt/lizardcd-server
EXPOSE 5117
ENTRYPOINT [ "./lizardcd-server" ]