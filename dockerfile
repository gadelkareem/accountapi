FROM golang:latest

RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go get -u ./...
RUN go build -o api .
RUN go test

ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.6.0/wait /wait
RUN chmod +x /wait

CMD /wait
