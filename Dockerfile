FROM golang:1.6.2-alpine
MAINTAINER zwh8800 <496781108@qq.com>

WORKDIR /app

RUN apk update && apk add ca-certificates && apk add git && \
    apk add tzdata && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone && go get github.com/Masterminds/glide

ADD ./template /app/template
ADD ./static /app/static
ADD . $GOPATH/src/github.com/zwh8800/md-blog-gen/

RUN cd $GOPATH/src/github.com/zwh8800/md-blog-gen && glide install && go install

VOLUME /app/log
VOLUME /app/config
VOLUME /app/static/img

CMD ["md-blog-gen", "-log_dir", "/app/log", "-alsologtostderr", "-config", "/app/config/golang-mirror.gcfg"]
