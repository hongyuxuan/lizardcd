FROM nginx:1.24

VOLUME /usr/share/nginx/html/lizardcd-doc/

COPY .vitepress/dist/ /usr/share/nginx/html/lizardcd-doc/

EXPOSE 80

RUN chmod -R 755 /usr/share/nginx/html/*