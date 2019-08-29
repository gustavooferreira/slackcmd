.PHONY: build
build:
	@go build ./...


.PHONY: test
test:
	@go test -v ./...


.PHONY: coverage
coverage:
	@go test -cover ./...


.PHONY: lint
lint:
	@go vet ./...


.PHONY: find_todo
find_todo:
	@grep --color=always --include=\*.go -PnRe '(//|/*).*TODO' ./ || true


.PHONY: find_fixme
find_fixme:
	@grep --color=always --include=\*.go -PnRe '(//|/*).*FIXME' ./ || true


.PHONY: find_xxx
find_xxx:
	@grep --color=always --include=\*.go -PnRe '(//|/*).*XXX' ./ || true


.PHONY: clean
clean:
	@# @rm -f file
	@echo "Removing files"


.PHONY: count
count:
	@echo "Lines of code:"
	@find . -type f -name "*.go" | xargs wc -l
