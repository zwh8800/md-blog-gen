FROM alpine:3.3
MAINTAINER zwh8800 <496781108@qq.com>

WORKDIR /app

RUN apk update && apk add ca-certificates && echo "Asia/Shanghai" > /etc/timezone

ADD ./md-blog-gen /app
ADD ./template /app/template
ADD ./static /app/static

VOLUME /app/log
VOLUME /app/config
VOLUME /app/static/img

EXPOSE 3336

CMD ["./md-blog-gen", "-log_dir", "log", "-config", "config/md-blog-gen.gcfg"]
