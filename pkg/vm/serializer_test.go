package vm

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSerialize(t *testing.T) {
	type args struct {
		name    string
		comment string
		b       []byte
	}
	tests := []struct {
		name    string
		args    args
		wantW   string
		wantErr bool
	}{
		{"test1", args{"test1", "foo", []byte{0, 1}}, "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			if err := Serialize(tt.args.name, tt.args.comment, tt.args.b, w); (err != nil) != tt.wantErr {
				t.Errorf("Serialize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			gotW := w.String()
			assert.Contains(t, gotW, `"name": "test1"`)
			assert.Contains(t, gotW, "AAE=")

			cb, err := Deserialize(strings.NewReader(gotW))
			assert.Nil(t, err)
			assert.Equal(t, "test1", cb.Name)
			assert.Equal(t, "foo", cb.Comment)
			assert.Equal(t, OpNop, cb.Data[0])
		})
	}
}
