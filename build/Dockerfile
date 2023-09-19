FROM golang:1.18.7 as builder

LABEL maintainer="max.liu9403@gmail.com"

ARG COMMITID
ENV COMMITID=${COMMITID:-v1.0}
ENV GO111MODULE=on
ENV MIRROR_URL=https://mirrors.aliyun.com/alpine/

RUN echo '' > /etc/apk/repositories \
    && echo "${MIRROR_URL}v3.14/main" >> /etc/apk/repositories \
    && echo "${MIRROR_URL}v3.14/community" >> /etc/apk/repositories

WORKDIR /go/src/github.com/maxliu9403/go-template

COPY . .

RUN go install -mod vendor -ldflags="-s -w -X 'main.Build=$COMMITID'" -v ./...

FROM golang:1.18.7

ENV MIRROR_URL=https://mirrors.aliyun.com/alpine/

RUN echo '' > /etc/apk/repositories \
    && echo "${MIRROR_URL}v3.14/main" >> /etc/apk/repositories \
    && echo "${MIRROR_URL}v3.14/community" >> /etc/apk/repositories

COPY --from=builder /go/bin/app /
COPY --from=builder /go/src/github.com/maxliu9403/go-template/docs/ /docs
COPY --from=builder /go/src/github.com/maxliu9403/go-template/configs/dev.yaml /

EXPOSE 9901

CMD ["/app", "--config", "dev.yaml"]
