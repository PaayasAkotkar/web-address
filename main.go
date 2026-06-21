package main

import (
	"app/webaddress/example"
	"log"
)

func main() {
	log.SetFlags(log.Lshortfile)
	example.PlayWebAddressURL()
	example.PlayWebAddress()

}
