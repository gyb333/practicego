package basic

import (
	"fmt"
	"reflect"
	"strings"
	. "unsafe"
)

func ReflectMain() {
	ret := CallMethod(func(arg interface{}) (string, bool) {
		fmt.Println(arg)
		return "test", true
	}, []interface{}{"gyb", 33})
	for i, v := range ret {
		fmt.Printf("%T,%v", v, v)
		fmt.Println(i, v,v.Interface(),v.Type())
		switch v.Interface().(type) {
		case string:
			fmt.Println(v.String(),reflect.TypeOf(""))
		case bool:
			fmt.Println(v.Bool(),reflect.TypeOf(false))
		}
	}

}

func ReflectData() {
	x := 1
	fmt.Printf("%T,%#v,%#x,%v\n", x, &x, Pointer(&x), x)
	d := reflect.ValueOf(&x).Elem() // d refers to the variable x
	fmt.Printf("%T,%#v,%#x,%v\n", d, &d, Pointer(&d), d)
	d.Set(reflect.ValueOf(2))
	fmt.Println(&x, x, d.Addr(), d)
	px := d.Addr().Interface().(*int) // px := &x
	*px = 3                           // x = 3
	fmt.Println(px, x, d)             // "3"
}

func Equal() {
	got := strings.Split("a:b:c", ":")
	want := []string{"a", "b", "c"}
	fmt.Println(reflect.DeepEqual(got, want))
	var a, b []string = nil, []string{}
	fmt.Println(reflect.DeepEqual(a, b)) // "false"

	var c, d map[string]int = nil, make(map[string]int)
	fmt.Println(reflect.DeepEqual(c, d)) // "false"

	fmt.Println(reflect.DeepEqual([]int{1, 2, 3}, []int{1, 2, 3}))        // "true"
	fmt.Println(reflect.DeepEqual([]string{"foo"}, []string{"bar"}))      // "false"
	fmt.Println(reflect.DeepEqual([]string(nil), []string{}))             // "false"
	fmt.Println(reflect.DeepEqual(map[string]int(nil), map[string]int{})) // "false"
}

func CallMethod(method interface{}, params interface{}) []reflect.Value {
	fv := reflect.ValueOf(method)
	args := []reflect.Value{reflect.ValueOf(params)}
	return fv.Call(args)
}
