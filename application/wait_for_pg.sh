#!/bin/bash
# 等待数据库服务
while ! nc -z db 5432; do
  sleep 1
done

# 数据库服务可用，启动应用
./app