FROM golang:1.16-alpine
WORKDIR /danmu/go-danmu
COPY . .

RUN go env -w GOPROXY=https://goproxy.cn,direct \
    && go mod tidy \
    && go build -o app main.go \
    && sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
    && apk update --no-cache \
    && apk add ffmpeg 

CMD ./app