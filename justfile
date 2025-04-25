# this is the default so that when you run `just` you get a nice interface
# for selecting a command
chose-a-command:
	just --choose

lint-nix-files:
	@echo "linting nix files"
	@out=0; \
		for i in $(find . -name '*.nix'); do \
			cat $i | nixfmt -c --filename "##$i" || out=$(( out + 1 ));  \
		done; \
		exit $out

lint: lint-nix-files
	golangci-lint run --config golangci.yaml

pre-pr: lint test

test:
	go test ./...

