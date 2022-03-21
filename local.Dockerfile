FROM golang:1.17

WORKDIR /code
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

EXPOSE 8080