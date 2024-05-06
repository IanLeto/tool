git commit -am "x"
export GOROOT=/home/ian/sdk/go1.22.2
GOOS=linux GOARCH=amd64 go build -tags netgo -ldflags '-extldflags "-static"' -o tool main.go && \
    docker build -t ianleto/tool:$(git rev-parse --short HEAD) -f dockerfile .&&\
    docker push ianleto/tool:$(git rev-parse --short HEAD)