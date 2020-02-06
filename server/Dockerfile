FROM golang:1.13.5-alpine3.10
ADD . /server/src
WORKDIR /server/src
RUN apk update
RUN apk add --no-cache git gcc libc-dev mercurial
RUN go get github.com/go-sql-driver/mysql github.com/jinzhu/gorm github.com/rs/cors
RUN cd /server/src && go build -o ../server
EXPOSE 5500
ENTRYPOINT /server/server