FROM alpine:3.3
MAINTAINER zwh8800 <496781108@qq.com>

WORKDIR /app
ADD ./md-blog-gen /app
ADD ./template /app/template

VOLUME /app/log
VOLUME /app/config

EXPOSE 3336

CMD ["./md-blog-gen", "-log_dir", "log", "-config", "config/md-blog-gen.gcfg"]
