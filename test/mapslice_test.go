package test

import (
	"fmt"
	"testing"
)

type data struct {
	num int
	key *string
	items map[string]bool

}
func (this *data) pointerFunc() {
	this.num = 7
}
func (this data) valueFunc() {
	this.num = 8
	*this.key = "valueFunc.key"
	this.items["valueFunc"] = true
}


func TestPointor(t *testing.T) {
	key := "key1"
	d := data{1, &key, make(map[string]bool)}
	fmt.Printf("num=%v key=%v items=%v\n", d.num, *d.key, d.items)
	d.pointerFunc()
	fmt.Printf("num=%v key=%v items=%v\n", d.num, *d.key, d.items)
	d.pointerFunc()
	fmt.Printf("num=%v key=%v items=%v\n", d.num, *d.key, d.items)
}