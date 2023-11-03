.PHONY: test-run
test-run:
	go build
	./terradiff -s branch1 -w ~/Desktop/terradiff-workspace -r 'https://github.com/paper2/test-terradiff' | jq -r '. | "\(.level): \(.msg)"'

