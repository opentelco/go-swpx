hash := $(shell git log --pretty=format:'%h' -n 1)
goarch := $(shell echo amd64)
goos :=  $(shell echo linux)
app_name := swpx
resources_bin := plugins/resources
providers_bin := plugins/providers
bin_path = ./bin
curpath := $(shell pwd)
resources_dir := $(curpath)/resources
providers_dir := $(curpath)/providers

resources_plugin_dir := $(curpath)/plugins/resources
providers_plugin_dir := $(curpath)/plugins/providers

.PHONY: build pb providers resources main clean linux hash

build: clean pb providers resources main


pb:
	# Generate new protobufs
	@go generate

providers:
	# Building PROVIDERS
	@rm -f $(providers_bin)/*
	@cd $(providers_dir)/vx/; go build -o $(providers_plugin_dir)/providers_vx .

resources:
	# Building RESOURCES
	@rm -f $(resources_bin)/*
	@cd $(resources_dir)/vrp_plugin/; go build -o $(resources_plugin_dir)/resource_vrp .
	@cd $(resources_dir)/raycore_plugin/; go build -o $(resources_plugin_dir)/resource_raycore .


main:
	# Building SWPX ($(bin_path)/$(app_name))
	@rm -f $(app_name)
	@go build -o $(bin_path)/$(app_name) main.go

clean_all:
	# Remove old binarys
	@rm -f $(resources_bin)/*
	@rm -f $(providers_bin)/*
	@rm -f $(app_name)

test:
	k6 run tests/api-req-k6.js -u 100 -d 5s



linux:

	rm -rf build/linux
	mkdir -p build/linux/plugins/resources;mkdir -p build/linux/plugins/providers;mkdir -p build/linux/config
	env GOOS=$(goos) GOARCH=$(goarch) go build -o ./build/linux/plugins/resources/resource_vrp_plug ./resources/vrp_plugin/
	env GOOS=$(goos) GOARCH=$(goarch) go build -o ./build/linux/plugins/providers/provider_vx ./providers/vx/main.go
	env GOOS=$(goos) GOARCH=$(goarch) go build -o ./build/linux/swpx main.go
	cp config/config.yml build/linux/config/
	zip -r build/swpx_$(goos)_$(goarch)-$(hash).zip ./build/linux/*
	@echo built and zipped build/swpx_$(goos)_$(goarch)-$(hash).zip



hash:
	@echo current hash is: $(hash)



