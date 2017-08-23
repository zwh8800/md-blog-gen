FROM golang:1.8-alpine
MAINTAINER zwh8800 <496781108@qq.com>

WORKDIR /app

RUN apk update && apk add ca-certificates && apk add git && \
    apk add tzdata && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone && go get -v github.com/Masterminds/glide

ADD ./template /app/template
ADD ./static /app/static
ADD . $GOPATH/src/github.com/zwh8800/md-blog-gen/

RUN cd $GOPATH/src/github.com/zwh8800/md-blog-gen && glide install && go install

VOLUME /app/log
VOLUME /app/config
VOLUME /app/static/img

EXPOSE 3336

CMD ["md-blog-gen", "-log_dir", "/app/log", "-alsologtostderr", "-config", "/app/config/md-blog-gen.toml"]
