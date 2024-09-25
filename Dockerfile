# 使用官方Go镜像作为构建环境
FROM golang:1.16-alpine AS builder
 
# 设置工作目录
WORKDIR /app
 
# 设置环境变量
ENV DB_USER=root \
DB_PASSWORD=123456 \
DB_NAME=demo \
DB_HOST=localhost \
DB_PORT=3306 \
# 阿里云访问配置
ALIBABA_CLOUD_ACCESS_KEY_ID=yourkey \
ALIBABA_CLOUD_ACCESS_KEY_SECRET=yoursecret \
Sign_Name=yourtemplatename \
Template_Code=yoursmscode \
# qq邮箱访问配置
SMTP_Username=837425169@qq.com \
SMTP_Password=yourauthcode \
 
# 复制go.mod和go.sum文件，以便于使用缓存
COPY go.mod go.sum ./
 
# 获取依赖项
RUN go mod download
 
# 复制项目源代码
COPY . .
 
# 编译构建应用程序
RUN go build -o /myapp .
 
# 使用scratch创建最小镜像
FROM scratch
 
# 复制构建阶段构建的应用程序
COPY --from=builder /myapp /myapp
 
# 暴露端口
EXPOSE 8080/tcp
 
# 设置容器启动时运行的命令
CMD ["/myapp"]