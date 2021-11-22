FROM golang:1.16-alpine
WORKDIR /danmu/go-danmu
COPY . .

RUN go env -w GOPROXY=https://goproxy.cn,direct

RUN go mod tidy
RUN go build -o app main.go

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

RUN apk update
RUN apk add ffmpeg

CMD ./app