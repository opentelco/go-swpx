hash := $(shell git log --pretty=format:'%h' -n 1)
GITTAG := $(shell git describe --tags --abbrev=0)
goarch := $(shell echo amd64)
goos :=  $(shell echo linux)
app_name := swpx
plugin_bin := plugins/resources
providers_bin := plugins/providers
bin_path = ./bin
curpath := $(shell pwd)
resources_dir := $(curpath)/resources
providers_dir := $(curpath)/providers

resources_plugin_dir := $(curpath)/plugins/resources
providers_plugin_dir := $(curpath)/plugins/providers

.PHONY: build pb providers resources core clean linux hash release docker

build: clean_all pb providers resources core


pb:
	# Generate new protobufs
	@go generate ./...

providers:
	# Building PROVIDERS
	@rm -f $(providers_bin)/*
	@cd $(providers_dir)/vx/; go build -o $(providers_plugin_dir)/vx .
	@cd $(providers_dir)/default/; go build -o $(providers_plugin_dir)/default .
	@cd $(providers_dir)/sait/; go build -o $(providers_plugin_dir)/sait .

plugins: r_clean r_vrp r_ctc r_generic

r_clean:
	@rm -f $(plugin_bin)/*

r_generic:
	# Building R-GENERIC
	@cd $(resources_dir)/generic/; go build -o $(resources_plugin_dir)/generic .

r_raycore:
	# Building R-RAYCORE
	@cd $(resources_dir)/raycore/; go build -o $(resources_plugin_dir)/raycore .

r_vrp:
	# Building R-VRP
	@cd $(resources_dir)/vrp/; go build -o $(resources_plugin_dir)/vrp .

r_ctc:
	# Building R-CTC
	@cd $(resources_dir)/ctc/; go build -o $(resources_plugin_dir)/ctc .

core:
	# Building SWPX ($(bin_path)/$(app_name))
	@rm -f $(app_name)
	@go build -o $(bin_path)/$(app_name) main.go

clean_all:
	# Remove old binarys
	@rm -f $(plugin_bin)/*
	@rm -f $(providers_bin)/*
	@rm -f $(app_name)

test:
	k6 run tests/api-req-k6.js -u 100 -d 5s



linux:
	rm -rf build/linux
	mkdir -p build/linux/plugins/resources;mkdir -p build/linux/plugins/providers;mkdir -p build/linux/config
	env GOOS=$(goos) GOARCH=$(goarch) go build -o ./build/linux/plugins/resources/vrp_plug ./resources/vrp/
	env GOOS=$(goos) GOARCH=$(goarch) go build -o ./build/linux/plugins/resources/ctc_plug ./resources/ctc/
	env GOOS=$(goos) GOARCH=$(goarch) go build -o ./build/linux/plugins/providers/provider_vx ./providers/vx/main.go
	env GOOS=$(goos) GOARCH=$(goarch) go build -o ./build/linux/swpx main.go
	cp config/config.yml build/linux/config/
	zip -r build/swpx_$(goos)_$(goarch)-$(hash).zip ./build/linux/*
	@echo built and zipped build/swpx_$(goos)_$(goarch)-$(hash).zip



hash:
	@echo current hash is: ${hash}

t:
	@echo current tag is: ${GITTAG}

docker:
	docker build -t registry.opentelco.io/go-swpx:${GITTAG} .
	docker push registry.opentelco.io/go-swpx:${GITTAG}

vxdocker:
	docker build -t 441617468760.dkr.ecr.eu-west-1.amazonaws.com/opentelco/swpx:${GITTAG} .
	docker push 441617468760.dkr.ecr.eu-west-1.amazonaws.com/opentelco/swpx:${GITTAG}

dockerhash:
	docker build -t registry.opentelco.io/go-swpx:${hash} .
	docker push registry.opentelco.io/go-swpx:${hash}
