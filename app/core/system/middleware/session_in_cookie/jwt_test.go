package session_in_cookie

import (
	"fmt"
	"testing"
)

func TestJWT_Decode(t *testing.T) {
	key := []byte("1234567890")
	a := PHPSessionId{15, 999999}
	b := a.Encode(key)
	fmt.Println(b)
	s := JWT("1kwsbd87xrtntypn3sdvnrkw7vo")
	c, err := s.Decode(key)
	fmt.Println(err)
	fmt.Println(c)
}

func TestPHPSessionId_Encode(t *testing.T) {
	type fields struct {
		UserIntId uint32
		StartTime uint32
	}
	type args struct {
		sharedKey []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &PHPSessionId{
				UserIntId: tt.fields.UserIntId,
				StartTime: tt.fields.StartTime,
			}
			if got := s.Encode(tt.args.sharedKey); got != tt.want {
				t.Errorf("Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}
