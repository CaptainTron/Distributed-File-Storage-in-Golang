package main

import (
	"bytes"
	"testing"
)

func TestCopyEncrypt(t *testing.T) {
	payload := "AES_ENCRYPTION"
	src := bytes.NewReader([]byte(payload))
	dst := new(bytes.Buffer)

	key := newEncryptKey()
	_, err := copyEncrypt(key, src, dst)
	if err != nil {
		t.Error(err)
	}

	out := new(bytes.Buffer)
	nw, err := copyDecrypt(key, dst, out)
	if err != nil {
		t.Error(err)
	}

	if nw != 16+len(payload) {
		t.Errorf("bytes do not match %d", nw)
	}

	if out.String() != payload {
		t.Errorf("decryption failed")
	}
}
