package main

import (
	"log"
	"xutils/xpredefine"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile | log.LUTC)
}

func main() {
	log.Println("xutils main start")
	xpredefine.SetDebug(true)
	log.Println("xutils main end")
}
