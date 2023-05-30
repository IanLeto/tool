# GoFrame Template For SingleRepo

Project Makefile Commands:

- `make cli`: Install or Update to the latest GoFrame CLI tool.
- `make dao`: Generate go files for `Entity/DAO/DO` according to the configuration file from `hack` folder.
- `make service`: Parse `logic` folder to generate interface go files into `service` folder.
- `make image TAG=xxx`: Run `docker build` to build image according `manifest/docker`.
- `make image.push TAG=xxx`: Run `docker build` and `docker push` to build and push image according `manifest/docker`.
- `make deploy TAG=xxx`: Run `kustomize build` to build and deploy deployment to kubernetes server group
  according `manifest/deploy`.

# docker build 本地

GOOS=linux GOARCH=amd64 go build -o ./ main.go
docker build --tag ianleto/goorigin:$(git rev-parse --short HEAD) -f ./manifest/docker/Dockerfile .
docker push ianleto/goorigin:$(git rev-parse --short HEAD)

# mock

mockgen -source internal/dao/es/es.go -destination testsuite/mock_apiClient.go -package="testsuite"