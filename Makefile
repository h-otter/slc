.PHONY: build
build:
	go build -o slc ./ && sudo setcap "CAP_SYS_ADMIN=ep" ./slc
