FROM golang:1.18.1-alpine

# 为我们的镜像设置必要的环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN apk add git
#添加ffmpeg
RUN apk update
RUN apk add yasm && apk add ffmpeg

# 移动到工作目录：/build
WORKDIR /build

# 复制项目中的 go.mod 和 go.sum文件并下载依赖信息
COPY go.mod .
COPY go.sum .
RUN go mod download

# 将代码复制到容器中
COPY . .

# 将我们的代码编译成二进制可执行文件 bubble
RUN go build -o douyin .



# 需要运行的命令
ENTRYPOINT ["/build/douyin"]