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

#### 使用 

#### 使用日志生成工具

    命令: ./bench file --content ian --rate 10 --size 10 
    说明: ./bench 使用fakelog生成日志 --content 日志内容 --rate 生成速率 --size 日志大小 path默认是标准输出

#### 使用 time

    命令: ./bench timeconv --key timestampToTime --value 1405000000
    说明: ./bench 使用time转化 --key 使用啥专访凡是,这个是时间戳转CST时间 --value 时间戳
    
    命令: ./bench timeconv --fomat 1 --value "2014-07-10 21:46:40 +0800 CST" 
    说明: ./bench 使用time转化 --format 使用啥模板,这个是CST模板 --target 时间格式

timeconv --key millisecondToTime --value 1690350196498 --opt add --params -1h