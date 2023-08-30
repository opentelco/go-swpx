package radius

import (
	"context"
	"log"
	"testing"

	"layeh.com/radius"
	"layeh.com/radius/rfc2865"
	"layeh.com/radius/rfc2866"
)

func Test_Server(t *testing.T) {
	go New()

	packet := radius.New(radius.CodeAccessRequest, []byte(`secret`))
	rfc2865.UserName_SetString(packet, "tim")
	rfc2865.UserPassword_SetString(packet, "12345")

	response, err := radius.Exchange(context.Background(), packet, "localhost:1812")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Code:", response.Code)

	packet = radius.New(radius.CodeAccountingRequest, []byte(`secret`))
	rfc2866.AcctAuthentic_Set(packet, rfc2866.AcctAuthentic_Value_RADIUS)
	response, err = radius.Exchange(context.Background(), packet, "localhost:1812")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Code:", response.Code)
}

func Test_Server_X(t *testing.T) {
	New()
}
