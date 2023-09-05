package radius

import (
	"fmt"
	"log"

	"git.liero.se/opentelco/go-swpx/fleet/aaa/radius/huawei"
	"layeh.com/radius"
	"layeh.com/radius/rfc2865"
)

func New() {

	handler := func(w radius.ResponseWriter, r *radius.Request) {

		switch r.Code {
		case radius.CodeAccessRequest:
			username := rfc2865.UserName_GetString(r.Packet)
			password := rfc2865.UserPassword_GetString(r.Packet)

			var code radius.Code
			fmt.Println(username, password)
			if username == "sven" && password == "bbnamnam" {
				code = radius.CodeAccessAccept

				rfc2865.ServiceType_Set(r.Packet, rfc2865.ServiceType_Value_AdministrativeUser)
				err := huawei.ExecPrivilege_Add(r.Packet, 15)
				if err != nil {
					log.Println(err)
				}

			} else {

				code = radius.CodeAccessReject
			}

			log.Printf("Writing %v to %v", code, r.RemoteAddr)
			w.Write(r.Response(code))

		case radius.CodeAccountingRequest:
			log.Println("Accounting-Request")

			w.Write(r.Response(radius.CodeAccountingResponse))

		default:
			log.Printf("Unknown code %v", r.Code)
		}

	}

	server := radius.PacketServer{
		Addr:         ":1813",
		Handler:      radius.HandlerFunc(handler),
		SecretSource: radius.StaticSecretSource([]byte(`secret`)),
	}

	log.Printf("Starting server on :1812")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
