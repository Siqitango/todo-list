#!/bin/bash

# 数据库配置
DB_USER="root"
DB_PASSWORD="123456"
DB_HOST="127.0.0.1"
DB_PORT="3306"
DB_NAME="todo_list"

# 执行SQL脚本
echo "正在初始化数据库..."
mysql -u ${DB_USER} -p${DB_PASSWORD} -h ${DB_HOST} -P ${DB_PORT} < $(dirname $0)/init.sql

# 检查执行结果
if [ $? -eq 0 ]; then
    echo "数据库初始化成功！"
else
    echo "数据库初始化失败，请检查配置和权限。"
    exit 1
fi