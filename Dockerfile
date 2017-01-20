FROM alpine
MAINTAINER Quey-Liang Kao <s101062801@m101.nthu.edu.tw>
# Last Modified: 2017/01/21

ENV PATH /bin:/sbin:/usr/bin

RUN apk update \
    && apk add go git

RUN mkdir -p /src/github.com/gammons/

RUN git clone https://github.com/gammons/todolist.git \
    && mv todolist /src/github.com/gammons/

ENV GOPATH /

RUN go install github.com/gammons/todolist

CMD ["todolist", "web"]

EXPOSE 7890


