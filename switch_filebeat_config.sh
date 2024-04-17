#!/bin/bash

# 设置变量
POD_NAME="your-pod-name"
CONTAINER_NAME="your-container-name"
LOCAL_FILE_PATH="/path/to/local/file"
REMOTE_FILE_PATH="/etc/filebeat/filebeat.yml"

# 使用 kubectl cp 命令复制文件
kubectl cp $LOCAL_FILE_PATH $POD_NAME:$REMOTE_FILE_PATH -c $CONTAINER_NAME
kubect exec -it $POD_NAME -c $CONTAINER_NAME -- /bin/bash
supervisord restart log