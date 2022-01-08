package identity

import (
	"app/utility/base32"
	"fmt"
	"reflect"
	"testing"
)

func TestObjectId(t *testing.T) {
	tests := []struct {
		name string
		want [12]byte
	}{
		// TODO: Add test cases.
	}
	for i := 0; i < 10; i++ {
		oid := ObjectId()
		fmt.Println(base32.EncodeId(oid[:]))
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ObjectId(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ObjectId() = %v, want %v", got, tt.want)
			}
		})
	}
}
