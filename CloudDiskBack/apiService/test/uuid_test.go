package test

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"testing"
)

func TestUuid(t *testing.T) {
	v4 := uuid.NewV4()
	fmt.Println(v4.String())
}
