package identity

import (
	"fmt"
	"net"
	"reflect"
	"testing"
)

func Test_mac(t *testing.T) {
	tests := []struct {
		name string
		want [3]byte
	}{
		// TODO: Add test cases.
	}
	interfaces, _ := net.Interfaces()
	for _, inter := range interfaces {

		fmt.Printf("%v %v %v\n", inter.Name, inter.HardwareAddr, inter.HardwareAddr.String())
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mac(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mac() = %v, want %v", got, tt.want)
			}
		})
	}
}
