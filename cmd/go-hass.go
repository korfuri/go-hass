package main

import (
	"fmt"
	"os"

	ghs "github.com/korfuri/go-hass/pkg/gohassapi"
)

func main() {
	hc := ghs.NewClient("https://hass.korfuri.fr/api/", os.Getenv("HASS_TOKEN"))
	status, err := hc.Check()
	if err != nil {
		fmt.Printf("API check error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(status)
	states, err := hc.States()
	if err != nil {
		fmt.Printf("API error while fetching states: %v\n", err)
		os.Exit(1)
	}
	for _, s := range states {
		fmt.Printf("State: %s is %s since %v.\n", s.EntityId, s.State, s.LastChanged)
	}
}
