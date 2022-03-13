package password

import (
	"fmt"
	"testing"
)

//func TestPasswordHash(t *testing.T) {
//	type args struct {
//		pwd     string
//		options *HashOptions
//	}
//	tests := []struct {
//		name string
//		args args
//		want string
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := Hash(tt.args.pwd, tt.args.options); got != tt.want {
//				t.Errorf("Hash() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestPasswordVerify(t *testing.T) {
//	type args struct {
//		pwd    string
//		hashed string
//	}
//	tests := []struct {
//		name    string
//		args    args
//		want    bool
//		wantErr bool
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, err := Verify(tt.args.pwd, tt.args.hashed)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("Verify() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if got != tt.want {
//				t.Errorf("Verify() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

func TestHash(t *testing.T) {
	password := "Admin@123456"
	for i := 0; i < 10; i++ {
		h := Hash(password, nil)
		fmt.Println(h)
		fmt.Println(Verify(password, h))
		fmt.Println(Verify("Admin#123456", h))
	}

}
