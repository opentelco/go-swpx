hash := $(shell git log --pretty=format:'%h' -n 1)
goarch := $(shell echo amd64)
goos :=  $(shell echo linux)
app_name := "swpx"
resources_bin := plugins/resources
providers_bin := plugins/providers

.PHONY: build linux hash

build:
	go generate
	rm -f $(resources_bin)/*
	rm -f $(providers_bin)/*
	rm -f $(app_name)
	# RESOURCES
	go build -o ./plugins/resources/resource_vrp_plug ./resources/vrp_plugin/
	go build -o ./plugins/resources/resource_raycore_plug ./resources/raycore_plugin/

	# PROVIDERS
	go build -o ./plugins/providers/provider_vx ./providers/vx/main.go
	

	# main
	go build -o ./bin/swpx main.go


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



