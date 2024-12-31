package core

import "encoding/hex"

type OID [12]byte

func (id OID) String() string {
	return string(id[:])
}

func (id OID) MarshalJSON() ([]byte, error) {
	return []byte(`"` + hex.EncodeToString(id[:]) + `"`), nil
}
