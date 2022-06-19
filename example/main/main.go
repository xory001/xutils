package main

import (
	"log"

	"github.com/xory001/xutils/xpredefine"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile | log.LUTC)
}

func main() {
	log.Println("xutils example start")
	xpredefine.SetDebug(true)
	log.Println("xutils example end")
}
