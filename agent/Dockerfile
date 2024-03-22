FROM alpine:3.15.4
ADD bin/lizardcd-agent /opt
WORKDIR /opt
RUN chmod +x /opt/lizardcd-agent
EXPOSE 5017
ENTRYPOINT [ "./lizardcd-agent" ]