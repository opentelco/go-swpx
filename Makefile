
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
	
	# go build -o ./plugins/providers/provider_default ./providers/default/main.go
	# go build -o ./plugins/providers/provider_ssab ./providers/ssab/main.go
	go build -o ./plugins/providers/provider_zitius ./providers/zitius/main.go
	
	# disabled
	# go build -o ./plugins/provider_telia_plug ./providers/telia/main.go
	# go build -o ./plugins/provider_ssab	_plug ./providers/ssab/main.go

	# main
	go build -o ./bin/swpx main.go


test:
	k6 run tests/api-req-k6.js -u 100 -d 5s
