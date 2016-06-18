package main

import "fmt"

func main() {
	api := phapi.New("", "")
	fmt.Println(api.Test())
}
