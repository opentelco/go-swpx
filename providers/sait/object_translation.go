package main

type Access struct {
	NetworkElement string
	Interface      string
	ResourcePlugin string
	SNMPCommunity  string
	Username       string
	Password       string
}

var translationMap = map[string]Access{
	"5f024967cf1c": Access{
		NetworkElement: "10.5.5.100",
		Interface:      "GigabitEthernet0/0/1",
		ResourcePlugin: "vrp_plugin",
		SNMPCommunity:  "xWTyZ9nA158ktJF2",
		Username:       "root",
		Password:       "qwerty1234",
	},
	"e2a4a1da985e": Access{
		NetworkElement: "10.5.5.100",
		Interface:      "GigabitEthernet0/0/2",
		SNMPCommunity:  "xWTyZ9nA158ktJF2",
		Username:       "root",
		Password:       "qwerty1234",
	},
	"83d636a256a6": Access{
		NetworkElement: "10.5.5.100",
		Interface:      "GigabitEthernet0/0/3",
		SNMPCommunity:  "xWTyZ9nA158ktJF2",
		Username:       "root",
		Password:       "qwerty1234",
	},
	"455ee346b701": Access{
		NetworkElement: "10.5.5.100",
		Interface:      "GigabitEthernet0/0/4",
		SNMPCommunity:  "xWTyZ9nA158ktJF2",
		Username:       "root",
		Password:       "qwerty1234",
	},
	"576655b914e7": Access{
		NetworkElement: "10.5.5.100",
		Interface:      "GigabitEthernet0/0/5",
		SNMPCommunity:  "xWTyZ9nA158ktJF2",
		Username:       "root",
		Password:       "qwerty1234",
	},
	"abc21bf112c6": Access{
		NetworkElement: "10.5.5.100",
		Interface:      "GigabitEthernet0/0/6",
		SNMPCommunity:  "xWTyZ9nA158ktJF2",
		Username:       "root",
		Password:       "qwerty1234",
	},
	"27c509a2c3a7": Access{
		NetworkElement: "10.5.5.100",
		Interface:      "GigabitEthernet0/0/7",
		SNMPCommunity:  "xWTyZ9nA158ktJF2",
		Username:       "root",
		Password:       "qwerty1234",
	},
	"df70884ef103": Access{
		NetworkElement: "10.5.5.100",
		Interface:      "GigabitEthernet0/0/8",
		SNMPCommunity:  "xWTyZ9nA158ktJF2",
		Username:       "root",
		Password:       "qwerty1234",
	},
	"2b9fc2048a8c": Access{
		NetworkElement: "10.5.5.100",
		Interface:      "GigabitEthernet0/0/9",
		SNMPCommunity:  "xWTyZ9nA158ktJF2",
		Username:       "root",
		Password:       "qwerty1234",
	},
	"abb77d79ad5d": Access{
		NetworkElement: "10.5.5.100",
		Interface:      "GigabitEthernet0/0/10",
		SNMPCommunity:  "xWTyZ9nA158ktJF2",
		Username:       "root",
		Password:       "qwerty1234",
	},
}
