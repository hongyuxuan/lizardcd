FROM alpine:3.19.1

RUN apk add sqlite

ADD bin/migrate /opt
ADD manifests /opt/manifests
WORKDIR /opt
RUN chmod +x /opt/migrate
ENTRYPOINT [ "./migrate" ]