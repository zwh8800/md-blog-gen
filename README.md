# Hello world

标签（空格分隔）： blog 闲扯

---

我把wordpress抛弃啦，旧的数据一会搞一个归档，url应该还不会变。

扔掉旧系统真是一身轻松啊！

老数据的链接：https://hzzz.lengzzz.com/

---

新的博客系统堪称相当轻量，花1天时间用golang写了个爬虫＋简单的博客系统。主要功能就是从cmd markdown上把我公开的文章爬过来存到本地。如果cmd markdown的作者能早日[开放API](https://github.com/ghosert/cmd-editor/issues/795)就好了。

另外话说cmd markdown真是宇宙最好的markdown编辑器，用着太舒服了，再次给作者点赞。

最后，这个爬虫我开源了，可能还会有希望基于markdown做博客系统的人。可以尝试一下。

```bash
go get -u -v github.com/zwh8800/md-blog-gen
md-blog-gen -log_dir log/ [-config <config.gcfg>]
open http://localhost:3336/
```

配置文件可以这么写
```ini
[dbConf]
driver="mysql"
dsn="username:password@tcp(mysql:3306)/mdblog?charset=utf8mb4&parseTime=true"

[env]
prod=true # 区分生产环境和测试环境
serverPort=3336

[spider]
startUrl="https://www.zybuluo.com/zwh8800/note/332154" # 你的随便一篇发布在cmd markdown的文章
spiderTag="blog" # 你要抓取的tag

[urlPush]
baidu="http://data.zz.baidu.com/urls?site=lengzzz.com&token=xxxxxxx" # 向百度搜索推送url收录的url

[site]
name="水能载舟 亦可赛艇" # 网站名称
description="冷zzZ和hzzZ的博客，热爱Linux，自称全栈，目前就职于美餐" # 描述
baseUrl="https://lengzzz.com/" # 网站首页
noteUrl="note" # 文字的url
tagUrl="tag" # tag的url
pageUrl="page" # page的url
rssUrl="rss"
staticUrl="static" # 静态文件的url
authorName="LengzzZ"
authorEmail="496781108@qq.com"
notePerPage=10
icp="津 ICP 备 14005803 号 - 2"

```
记着把注释删掉哦

另外最近用docker比较多，我正在把我所有用到的东西向docker上迁。所以这个爬虫也写了Dockerfile。你用docker的话可以这样部署。
```bash
git clone https://github.com/zwh8800/md-blog-gen
cd md-blog-gen/
CGO_ENABLED=0 go build
docker build -t zwh8800/md-blog-gen .
docker run -d -v log:/app/log -v config:/app/config --name blog zwh8800/md-blog-gen
```
图省事可以从docker hub上拉
```bash
docker pull zwh8800/md-blog-gen/
# 国内可以拉这个 速度快一些
docker pull registry.aliyuncs.com/zwh8800/md-blog-gen
docker run -d -v log:/app/log -v config:/app/config --name blog zwh8800/md-blog-gen

```

服务启停使用：
```bash
docker start blog
docker stop blog
docker restart blog
```

---

另外，最近换上了http2和let's encrypt的https证书。一个字，爽。
基本参照一位[大神的博客](https://imququ.com/)搞得。另外这个博客里干货很多，值得一看。
