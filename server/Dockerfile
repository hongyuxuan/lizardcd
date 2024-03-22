FROM alpine:3.15.4
ADD bin/lizardcd-server /opt
WORKDIR /opt
RUN chmod +x /opt/lizardcd-server
EXPOSE 5117
ENTRYPOINT [ "./lizardcd-server" ]