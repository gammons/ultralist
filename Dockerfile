FROM golang:alpine

RUN apk update && apk add git

RUN mkdir -p /src/github.com/gammons/

RUN go-wrapper download github.com/gammons/todolist
RUN go-wrapper install github.com/gammons/todolist

RUN ln -s /.todos.json $HOME/.todos.json

CMD ["todolist", "web"]

EXPOSE 7890
