FROM golang:latest
RUN apt-get update 
WORKDIR /telebot
ADD . /telebot
RUN make build
ENTRYPOINT ["bin/telebot"]
