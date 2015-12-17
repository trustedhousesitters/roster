FROM golang:1.5.2
COPY /app ${GOPATH}/src/app
RUN go get github.com/Masterminds/glide
RUN go build github.com/Masterminds/glide
WORKDIR ${GOPATH}/src/app
RUN GO15VENDOREXPERIMENT=1 glide install
RUN GO15VENDOREXPERIMENT=1 go build

EXPOSE 3333

ENTRYPOINT ["./app"]
