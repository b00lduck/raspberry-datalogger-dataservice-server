FROM hypriot/rpi-golang
WORKDIR /gopath1.5/src/github.com/b00lduck/raspberry-datalogger-dataservice-server
CMD ["raspberry-datalogger-dataservice-server"]

EXPOSE 8080
ADD . /gopath1.5/src/github.com/b00lduck/raspberry-datalogger-dataservice-server/
RUN go get
RUN go build
