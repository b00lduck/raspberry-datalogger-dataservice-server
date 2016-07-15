FROM hypriot/rpi-golang
WORKDIR /gopath1.5/src/b00lduck/datalogger/dataservice
CMD ["dataservice"]

EXPOSE 8080
ADD . /gopath1.5/src/b00lduck/datalogger/dataservice
RUN go get
RUN go build
