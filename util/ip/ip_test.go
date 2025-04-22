package ip

import (
	"fmt"
	"testing"
)

func TestGet(t *testing.T) {
	ip, err := Get()
	if err != nil {
		panic(err)
	}
	fmt.Println(ip)
}

func TestQuery(t *testing.T) {
	info, err := Query("1.1.1.1")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", info)
}
