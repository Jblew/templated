FROM golang:1.17-alpine
WORKDIR /app
ADD . /app

RUN GOBIN=/bin/ CGO_ENABLED=0 go install

CMD ["templated"]