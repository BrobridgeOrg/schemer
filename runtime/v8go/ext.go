package v8go_runtime

import (
	"bytes"

	msgpack "github.com/vmihailenco/msgpack/v5"
)

/*
type Uint8Array []byte

	func (u Uint8Array) MarshalMsgpack() ([]byte, error) {
		return []byte(u), nil
	}

	func (u *Uint8Array) UnmarshalMsgpack(b []byte) error {
		*u = append((*u)[:0], b...)
		return nil
	}
*/
type Uint8Array struct {
	Buffer bytes.Buffer
}

func (u Uint8Array) MarshalMsgpack() ([]byte, error) {
	return u.Buffer.Bytes(), nil
}

func (u *Uint8Array) UnmarshalMsgpack(b []byte) error {
	u.Buffer.Reset()
	_, err := u.Buffer.Write(b)
	return err
}

func init() {
	//msgpack.RegisterExt(18, (*Uint8Array)(nil))
	msgpack.RegisterExt(18, new(Uint8Array))
}
