.PHONY: run
run:
	go build
	./teradiff --branch branch1 | jq -r '. | "\(.level): \(.msg)"'

