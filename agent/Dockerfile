FROM alpine:3.19.1
ADD bin/lizardcd-agent /opt
WORKDIR /opt
RUN chmod +x /opt/lizardcd-agent
EXPOSE 5017
ENTRYPOINT [ "./lizardcd-agent" ]