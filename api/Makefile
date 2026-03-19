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
	jq '.abi' $(OUT)/AlphaToken.sol/AlphaToken.json > /tmp/AlphaToken_abi.json
	$(ABIGEN) --abi /tmp/AlphaToken_abi.json --pkg bindings --type AlphaToken --out internal/chain/bindings/alpha_token.go
	jq '.abi' $(OUT)/SubnetNFT.sol/SubnetNFT.json > /tmp/SubnetNFT_abi.json
	$(ABIGEN) --abi /tmp/SubnetNFT_abi.json --pkg bindings --type SubnetNFT --out internal/chain/bindings/subnet_n_f_t.go
	jq '.abi' $(OUT)/StakingVault.sol/StakingVault.json > /tmp/StakingVault_abi.json
	$(ABIGEN) --abi /tmp/StakingVault_abi.json --pkg bindings --type StakingVault --out internal/chain/bindings/staking_vault.go
	jq '.abi' $(OUT)/AWPEmission.sol/AWPEmission.json > /tmp/AWPEmission_abi.json
	$(ABIGEN) --abi /tmp/AWPEmission_abi.json --pkg bindings --type AWPEmission --out internal/chain/bindings/a_w_p_emission.go
	jq '.abi' $(OUT)/AWPDAO.sol/AWPDAO.json > /tmp/AWPDAO_abi.json
	$(ABIGEN) --abi /tmp/AWPDAO_abi.json --pkg bindings --type AWPDAO --out internal/chain/bindings/a_w_p_d_a_o.go
	jq '.abi' $(OUT)/StakeNFT.sol/StakeNFT.json > /tmp/StakeNFT_abi.json
	$(ABIGEN) --abi /tmp/StakeNFT_abi.json --pkg bindings --type StakeNFT --out internal/chain/bindings/stake_n_f_t.go

generate-sqlc:
	$(SQLC) generate

migrate:
	atlas schema apply --env local

build:
	go build -o bin/api ./cmd/api
	go build -o bin/indexer ./cmd/indexer
	go build -o bin/keeper ./cmd/keeper

test:
	go test ./internal/chain/ -count=1 -timeout 60s
	go test ./internal/server/handler/ -count=1 -timeout 60s
	go test ./internal/server/ws/ -count=1 -timeout 30s

lint:
	$(shell go env GOPATH)/bin/golangci-lint run ./...

clean:
	rm -rf bin/
