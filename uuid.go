package main

import (
	"github.com/gofrs/uuid"
)

func uuidV5Gen(ns int, h string) string {
	var n uuid.UUID

	switch ns {
	// DNS/Hostname
	case 1:
		n, _ = uuid.FromString("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	// URL
	case 2:
		n, _ = uuid.FromString("6ba7b811-9dad-11d1-80b4-00c04fd430c8")
	// OID
	case 3:
		n, _ = uuid.FromString("6ba7b812-9dad-11d1-80b4-00c04fd430c8")
	// Custom/X500
	default:
		n, _ = uuid.FromString("6ba7b814-9dad-11d1-80b4-00c04fd430c8")
	}

	u5 := uuid.NewV5(n, h)

	return u5.String()
}
