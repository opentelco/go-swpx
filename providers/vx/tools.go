package main

import "strings"

func HostFromFQDN(fqdn string) string {
	x := strings.Split(fqdn, ".")
	if len(x) > 0 {
		return x[0]
	}
	return fqdn
}
