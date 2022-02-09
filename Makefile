all:
	go build -ldflags "-s -w" -o relayer identity-token-relayer
