.PHONY: generate migrate build test lint clean

FORGE := $(HOME)/.foundry/bin/forge
ABIGEN := $(shell go env GOPATH)/bin/abigen
SQLC := $(shell go env GOPATH)/bin/sqlc
CONTRACTS := ../contracts
OUT := $(CONTRACTS)/out

generate: generate-bindings generate-sqlc

generate-bindings:
	cd $(CONTRACTS) && $(FORGE) build
	@mkdir -p internal/chain/bindings
	jq '.abi' $(OUT)/AWPRegistry.sol/AWPRegistry.json > /tmp/AWPRegistry_abi.json
	$(ABIGEN) --abi /tmp/AWPRegistry_abi.json --pkg bindings --type AWPRegistry --out internal/chain/bindings/awp_registry.go
	jq '.abi' $(OUT)/AWPToken.sol/AWPToken.json > /tmp/AWPToken_abi.json
	$(ABIGEN) --abi /tmp/AWPToken_abi.json --pkg bindings --type AWPToken --out internal/chain/bindings/a_w_p_token.go
	jq '.abi' $(OUT)/WorknetToken.sol/WorknetToken.json > /tmp/WorknetToken_abi.json
	$(ABIGEN) --abi /tmp/WorknetToken_abi.json --pkg bindings --type WorknetToken --out internal/chain/bindings/worknet_token.go
	jq '.abi' $(OUT)/AWPWorkNet.sol/AWPWorkNet.json > /tmp/AWPWorkNet_abi.json
	$(ABIGEN) --abi /tmp/AWPWorkNet_abi.json --pkg bindings --type AWPWorkNet --out internal/chain/bindings/awp_work_net.go
	jq '.abi' $(OUT)/AWPAllocator.sol/AWPAllocator.json > /tmp/AWPAllocator_abi.json
	$(ABIGEN) --abi /tmp/AWPAllocator_abi.json --pkg bindings --type AWPAllocator --out internal/chain/bindings/awp_allocator.go
	jq '.abi' $(OUT)/AWPEmission.sol/AWPEmission.json > /tmp/AWPEmission_abi.json
	$(ABIGEN) --abi /tmp/AWPEmission_abi.json --pkg bindings --type AWPEmission --out internal/chain/bindings/a_w_p_emission.go
	jq '.abi' $(OUT)/AWPDAO.sol/AWPDAO.json > /tmp/AWPDAO_abi.json
	$(ABIGEN) --abi /tmp/AWPDAO_abi.json --pkg bindings --type AWPDAO --out internal/chain/bindings/a_w_p_d_a_o.go
	jq '.abi' $(OUT)/veAWP.sol/veAWP.json > /tmp/veAWP_abi.json
	$(ABIGEN) --abi /tmp/veAWP_abi.json --pkg bindings --type VeAWP --out internal/chain/bindings/ve_awp.go

generate-sqlc:
	$(SQLC) generate

migrate:
	atlas schema apply --env local

VERSION ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X github.com/cortexia/rootnet/api/internal/server/handler.Version=$(VERSION)"

build:
	go build $(LDFLAGS) -o bin/api ./cmd/api
	go build $(LDFLAGS) -o bin/indexer ./cmd/indexer
	go build $(LDFLAGS) -o bin/keeper ./cmd/keeper

test:
	go test ./internal/chain/ -count=1 -timeout 60s
	go test ./internal/server/handler/ -count=1 -timeout 60s
	go test ./internal/server/ws/ -count=1 -timeout 30s

lint:
	$(shell go env GOPATH)/bin/golangci-lint run ./...

clean:
	rm -rf bin/
