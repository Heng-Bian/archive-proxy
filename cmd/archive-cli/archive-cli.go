package main

import "fmt"

func main() {
	s:=make([]string,1)
	s=nil
	fmt.Printf("%v",s)
	fmt.Printf("%v",s == nil)
	fmt.Printf("%v",len(s))
}
