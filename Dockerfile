FROM alpine
MAINTAINER Bumwoo Park <qjadn0914@naver.com>

RUN apk update && apk add --no-cache git
RUN mkdir /images
ADD ./sumwhere /usr/local/bin
ENTRYPOINT ["/usr/local/bin/sumwhere"]

