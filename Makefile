BUILD_DIR=build
EXEC_NAME_PREFIX=firmware_go
BULD_FLAGS=-tags "fts5"
LDLFAGS="-s -w"
CGO_ENABLED=0

DARWIN_PACK=$(EXEC_NAME_PREFIX)_darwin_amd64
LINUX_PACK=$(EXEC_NAME_PREFIX)_linux_amd64
WIN_PACK=$(EXEC_NAME_PREFIX)_windows_amd64.exe

.PHONY: all


buildAll: clean ${BUILD_DIR}/app

${BUILD_DIR}/app:
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=$(CGO_ENABLED) go build $(BULD_FLAGS) -ldflags $(LDLFAGS) -o $(BUILD_DIR)/$(DARWIN_PACK) ./main.go
	#GOOS=linux GOARCH=amd64 CGO_ENABLED=$(CGO_ENABLED) go build $(BULD_FLAGS) -ldflags $(LDLFAGS) -o $(LINUX_PACK)
	#GOOS=windows GOARCH=amd64 CGO_ENABLED=$(CGO_ENABLED) go build $(BULD_FLAGS) -ldflags $(LDLFAGS) -o $(WIN_PACK)

clean:
	rm -rf build/*
