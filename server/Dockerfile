FROM repo.cicc.com.cn/fi-public-docker-local/alpine:3.15.4-20220527
ADD bin/lizardcd-server /opt
WORKDIR /opt
RUN chmod +x /opt/lizardcd-server
EXPOSE 5117
ENTRYPOINT [ "./lizardcd-server" ]