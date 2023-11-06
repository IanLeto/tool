git commit -am "x"
GOOS=linux GOARCH=amd64 go build -o ori main.go && \
    docker build -t ianleto/tool:$(git rev-parse --short HEAD) -f Dockerfile2 .&&\
    docker push ianleto/tool:$(git rev-parse --short HEAD)