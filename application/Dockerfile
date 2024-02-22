FROM golang:1.20 AS app
ENV GO111MODULE=on
ENV GOPROXY https://goproxy.cn,direct
WORKDIR /root/togettoyou
COPY server/. .
RUN CGO_ENABLED=0 go build -v -o "app" .

FROM node:12 AS web
WORKDIR /root/togettoyou
COPY web/package*.json ./
RUN yarn config set registry https://registry.npmmirror.com/ && yarn install
COPY web/. .
RUN yarn run build:prod

FROM alpine:latest AS final
WORKDIR /root/togettoyou/

RUN apk add --no-cache netcat-openbsd iputils

COPY --from=app /root/togettoyou/app ./
COPY --from=app /root/togettoyou/config.yaml ./
COPY --from=web /root/togettoyou/dist/ ./dist/

ENTRYPOINT ["./app"]