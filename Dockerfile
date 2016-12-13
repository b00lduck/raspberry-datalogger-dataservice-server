FROM rem/rpi-golang-1.7:latest

WORKDIR /gopath/src/github.com/b00lduck/raspberry-datalogger-dataservice-server
ENTRYPOINT ["raspberry-datalogger-dataservice-server"]

EXPOSE 8080
ADD . /gopath/src/github.com/b00lduck/raspberry-datalogger-dataservice-server/
RUN go get
RUN go build

