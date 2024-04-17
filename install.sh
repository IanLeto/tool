#!/bin/bash

# 检查是否提供了参数
if [ "$#" -ne 2 ]; then
    echo "You must provide exactly two arguments: the directory name and an environment variable key=value pair."
    exit 1
fi

# 获取参数
dir=$1
env_var=$2

# 拆分环境变量键值对
IFS='=' read -ra KV <<< "$env_var"

# 检查键值对是否有效
if [ "${#KV[@]}" -ne 2 ]; then
    echo "Invalid environment variable. Must be in the format key=value"
    exit 1
fi

# 设置环境变量
export ${KV[0]}=${KV[1]}

# 检查目录是否存在
if [ ! -d "./$dir" ]; then
    echo "Directory ./$dir does not exist"
    exit 1
fi

# 检查 configmap.yaml 文件是否存在
if [ ! -f "./$dir/configmap.yaml" ]; then
    echo "File ./$dir/configmap.yaml does not exist"
    exit 1
fi

# 应用 configmap.yaml 文件
kubectl apply -f ./$dir/configmap.yaml
#env
echo $cluster