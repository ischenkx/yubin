FROM golang:latest

EXPOSE 6060

ADD bin ./bin
ADD run.sh ./run.sh
ADD config.yml ./config.yml

CMD bash ./run.sh config.yml