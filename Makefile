
build:
	go generate
	rm -f plugins/resources/*
	# rm -f plugins/providers/*
	# RESOURCES
	go build -o ./plugins/resources/resource_vrp_plug ./resources/vrp_plugin/
	# disabled
	# go build -o ./plugins/resource_raycore_plug ./resources/raycore_plugin/main.go
	# go build -o ./plugins/resource_comware_plug ./resources/comware_plugin/main.go
	
	# PROVIDERS
	go build -o ./plugins/providers/provider_vx ./providers/vx/main.go
	

	# main
	go build -o ./bin/swpx main.go


test:
	k6 run tests/api-req-k6.js -u 100 -d 5s



linux:
	rm -rf build/linux
	mkdir -p build/linux/plugins/resources
	mkdir -p build/linux/plugins/providers
	env GOOS=linux GOARCH=amd64 go build -o ./build/linux/plugins/resources/resource_vrp_plug ./resources/vrp_plugin/
	env GOOS=linux GOARCH=amd64 go build -o ./build/linux/plugins/providers/provider_x ./providers/vx/main.go
	env GOOS=linux GOARCH=amd64 go build -o ./build/linux/swpx main.go


