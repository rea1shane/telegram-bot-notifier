package ip

import (
	"fmt"
	"testing"
)

func TestGet(t *testing.T) {
	info, err := Get()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", info)
}
