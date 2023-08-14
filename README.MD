#### ianleto 测试用工具

go build -o ./bench main.go
#### 编译
GOOS=linux GOARCH=amd64 go build -o ./bench main.go && \
docker build --tag ianleto/bench:$(git rev-parse --short HEAD) -f ./dockerfile . && \
docker push ianleto/bench:$(git rev-parse --short HEAD) && \
k delete -f deployment.yaml && \
k apply -f deployment.yaml

#### 使用说明
json 去空格:
./bench jsonconv -r -i '\
{
"policy": {
"phases": {
"hot": {
"min_age": "0ms",
"actions": {
"rollover": {
"max_age": "1d",
"max_size": "50gb"
},
"set_priority": {
"priority": 100
}
}
}
}
}
}'


