package main

import (
	"log"

	"github.com/gofrs/uuid"
)

func main() {
	var ns uuid.UUID

	namespace := "URL"
	name := "test"

	switch namespace {
	case "DNS":
		ns, _ = uuid.FromString("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	case "URL":
		ns, _ = uuid.FromString("6ba7b811-9dad-11d1-80b4-00c04fd430c8")
	case "OID":
		ns, _ = uuid.FromString("6ba7b812-9dad-11d1-80b4-00c04fd430c8")
	case "X500":
		ns, _ = uuid.FromString("6ba7b814-9dad-11d1-80b4-00c04fd430c8")
	}


	u3 := uuid.NewV3(ns, name)
	log.Printf("UUIDv5: %s\n", u3.String())
	log.Printf("UUIDv5: %d\n", u3.Version())

	u5 := uuid.NewV5(ns, name)
	log.Printf("UUIDv5: %s\n", u5.String())
	log.Printf("UUIDv5: %d\n", u5.Version())

	startServer()
}
