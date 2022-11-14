package crypto

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCrypto_Encrypt(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name string
		c    *Crypto
		args args
	}{
		{
			name: "pos",
			c:    NewCrypto([]byte("12345678123456781234567812345678")),
			args: args{
				data: []byte("plain text"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ciphered := tt.c.Encrypt(tt.args.data)
			require.NotEqual(t, ciphered, tt.args.data)
			plain := tt.c.Decrypt(ciphered)
			require.Equal(t, plain, tt.args.data)

		})
	}
}
