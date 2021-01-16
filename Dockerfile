FROM golang:1.15.6 as builder

ARG GOPROXY=https://goproxy.cn
ENV GOPROXY $GOPROXY
RUN mkdir -p /root/workspace
WORKDIR /root/workspace

# restore dependencies
COPY . .
RUN go build -o app
RUN find . \( -name app -or -name "config*" \) | xargs tar cvfz app.tar.gz

FROM hidevops/base-go:v1.0.0 as release
COPY --from=builder /root/workspace/app.tar.gz /opt/app-root
RUN tar xvf /opt/app-root/app.tar.gz

EXPOSE 8080 7575
ENTRYPOINT ["/opt/app-root/app"]
