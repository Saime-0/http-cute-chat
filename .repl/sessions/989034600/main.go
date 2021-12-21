package main

import (
	"fmt"
	"time"
)

var ()

func main() {
	fmt.Printf("<%T> %+v\n", int64(time.Hour*24*10/time.Second), int64(time.Hour*24*10/time.Second))
}
