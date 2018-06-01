FROM alpine:latest

MAINTAINER huangjia huangjie@ennew.com

WORKDIR /

RUN ln -snf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

ADD ./nfs-util /
ADD ./start.sh /

ENTRYPOINT ["sh","-c","/start.sh"]

