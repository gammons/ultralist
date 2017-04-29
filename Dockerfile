FROM golang:alpine
MAINTAINER Quey-Liang Kao <s101062801@m101.nthu.edu.tw>
# Last Modified: 2017/01/21

RUN apk update && apk add git

RUN mkdir -p /src/github.com/gammons/

RUN go-wrapper download github.com/gammons/todolist
RUN go-wrapper install github.com/gammons/todolist

RUN ln -s /.todos.json $HOME/.todos.json

CMD ["todolist", "web"]

EXPOSE 7890
