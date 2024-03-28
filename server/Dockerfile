FROM repo.cicc.com.cn/public-docker-virtual/alpine:3.19.1

RUN set -ex \
        && alpine_ver=$(awk '/Linux/{print $NF}' /etc/issue) \
        && echo "https://repo.cicc.com.cn/artifactory/public-org-alpine-remote/v${alpine_ver}/main" >/etc/apk/repositories \
        && echo "https://repo.cicc.com.cn/artifactory/public-org-alpine-remote/v${alpine_ver}/community" >>/etc/apk/repositories \
	      && apk add sqlite
        
ADD bin/lizardcd-server /opt
WORKDIR /opt
RUN chmod +x /opt/lizardcd-server
EXPOSE 5117
ENTRYPOINT [ "./lizardcd-server" ]