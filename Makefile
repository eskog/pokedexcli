build:
	@go build -o bin/pokedexcli;

run: build
	@./bin/pokedexcli;

test:
	@go test ./... -v;
