package main

import (
	"fmt"

	"github.com/PlanetHoster/api_Golang/phapi"
)

func main() {
	api := phapi.New("", "")
	fmt.Println(api.Test())
}
