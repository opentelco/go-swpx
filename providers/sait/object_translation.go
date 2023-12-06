package main

type Access struct {
	Device         string
	Interface      string
	ResourcePlugin string
	SNMPCommunity  string
	Username       string
	Password       string
}

var translationMap = map[string]Access{
	"5f024967cf1c": Access{
		Device:         "10.5.5.100",
		Interface:      "GigabitEthernet0/0/1",
		ResourcePlugin: "vrp",
		Username:       "root",
		Password:       "qwerty1234",
	},
	"e2a4a1da985e": Access{
		Device:         "10.5.5.100",
		Interface:      "GigabitEthernet0/0/2",
		ResourcePlugin: "vrp",
		Username:       "root",
		Password:       "qwerty1234",
	},
	"83d636a256a6": Access{
		Device:         "10.5.5.100",
		Interface:      "GigabitEthernet0/0/3",
		ResourcePlugin: "vrp",

		Username: "root",
		Password: "qwerty1234",
	},
	"455ee346b701": Access{
		Device:         "10.5.5.100",
		Interface:      "GigabitEthernet0/0/4",
		ResourcePlugin: "vrp",

		Username: "root",
		Password: "qwerty1234",
	},
	"576655b914e7": Access{
		Device:         "10.5.5.100",
		Interface:      "GigabitEthernet0/0/5",
		ResourcePlugin: "vrp",

		Username: "root",
		Password: "qwerty1234",
	},
	"abc21bf112c6": Access{
		Device:         "10.5.5.100",
		Interface:      "GigabitEthernet0/0/6",
		ResourcePlugin: "vrp",

		Username: "root",
		Password: "qwerty1234",
	},
	"27c509a2c3a7": Access{
		Device:         "10.5.5.100",
		Interface:      "GigabitEthernet0/0/7",
		ResourcePlugin: "vrp",

		Username: "root",
		Password: "qwerty1234",
	},
	"df70884ef103": Access{
		Device:         "10.5.5.100",
		Interface:      "GigabitEthernet0/0/8",
		ResourcePlugin: "vrp",

		Username: "root",
		Password: "qwerty1234",
	},
	"2b9fc2048a8c": Access{
		Device:         "10.5.5.100",
		Interface:      "GigabitEthernet0/0/9",
		ResourcePlugin: "vrp",

		Username: "root",
		Password: "qwerty1234",
	},
	"abb77d79ad5d": Access{
		Device:         "10.5.5.100",
		Interface:      "GigabitEthernet0/0/10",
		ResourcePlugin: "vrp",

		Username: "root",
		Password: "qwerty1234",
	},
}
