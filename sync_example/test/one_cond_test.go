package test_test

import (
	. ".."
	"fmt"
	"testing"
)


func TestCondition(t *testing.T) {
	oneCond := NewOneCondition()
	go func() {
		for i:=1;i<=50;i++ {
			oneCond.DoNextFunc(func() {
				for j:=1;j<=1;j++ {
					fmt.Printf("Next go sequence of %d,loop of %d\n" ,j, i);
				}
			})
			//oneCond.DoFunc(true,func() {
			//	for j:=1;j<=3;j++ {
			//		fmt.Printf("Next go sequence of %d,loop of %d\n" ,j, i);
			//	}
			//})
		}
	}()
	for i:=1;i<=50;i++ {
		oneCond.DoFrontFunc(func() {
			for j:=1;j<=1;j++ {
				fmt.Printf("Front go sequence of %d,loop of %d\n" ,j, i);
			}
		})
		//oneCond.DoFunc(false,func() {
		//	for j:=1;j<=5;j++ {
		//		fmt.Printf("Front go sequence of %d,loop of %d\n" ,j, i);
		//	}
		//})
	}
}
