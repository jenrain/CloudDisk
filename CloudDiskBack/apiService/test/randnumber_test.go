package test

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestRandNumber(t *testing.T) {
	fmt.Println(rand.Perm(5))
}
