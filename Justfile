

install:
    go install -v ./cmd/oc/...

test:
    go test -v ./...

lint:
    go fmt ./...
    go vet ./...

api:
    curl https://api.stoplight.io/projects/cHJqOjI3NzgwNQ/branches/main/export/reference/exchange.yaml > docs/exchange.yaml
    oapi-codegen -config ./oapi-codegen.yaml -package api docs/exchange.yaml | grep -v WARNING > server/client/api/api.gen.go
