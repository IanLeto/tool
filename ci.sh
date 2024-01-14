git commit -am "x"
GOOS=linux GOARCH=amd64 go build -o tool main.go && \
    docker build -t ianleto/tool:$(git rev-parse --short HEAD) -f dockerfile .&&\
    docker push ianleto/tool:$(git rev-parse --short HEAD)