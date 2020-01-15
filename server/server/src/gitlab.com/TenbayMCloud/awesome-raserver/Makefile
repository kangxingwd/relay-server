GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=awesome-raserver
INSTALL_DST=$(GOPATH)/bin/
CONFIG_FILE=conf/cfg.toml

all: build
build:
	$(GOBUILD) -o $(BINARY_NAME) -v
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
install:
	cp $(BINARY_NAME) $(INSTALL_DST)
	cp $(CONFIG_FILE) $(INSTALL_DST)

run:
	$(GOBUILD) -o $(BINARY_NAME) -v
	./$(BINARY_NAME)
deps:
	$(GOGET) github.com/cihub/seelog
	$(GOGET) github.com/fasthttp/router
	$(GOGET) github.com/jinzhu/gorm
	$(GOGET) github.com/mitchellh/mapstructure
	$(GOGET) github.com/spf13/viper
	$(GOGET) github.com/valyala/fasthttp
	$(GOGET) github.com/go-sql-driver/mysql
	$(GOGET) github.com/go-redis/redis
	$(GOGET) github.com/json-iterator/go
