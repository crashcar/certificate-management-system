FROM golang:1.20 AS app
ENV GO111MODULE=on
ENV GOPROXY https://goproxy.cn,direct
WORKDIR /root/chiikawa
COPY server/. .
RUN CGO_ENABLED=0 go build -v -o "app" .

FROM node:12 AS web
WORKDIR /root/chiikawa
COPY web/package*.json ./
RUN yarn config set registry https://registry.npmmirror.com/ && yarn install
COPY web/. .
RUN yarn run build:prod

FROM alpine:latest AS final
WORKDIR /root/chiikawa/

RUN apk add --no-cache netcat-openbsd iputils

COPY --from=app /root/chiikawa/app ./
COPY --from=app /root/chiikawa/aesKey.txt ./
COPY --from=app /root/chiikawa/config.yaml ./
COPY --from=web /root/chiikawa/dist/ ./dist/

ENTRYPOINT ["./app"]
