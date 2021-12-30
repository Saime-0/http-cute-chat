package main

import (
	"fmt"
	"regexp"
)

var (
	re = regexp.MustCompile(`(?m)(0-9a-fA-F3)1,2`)
)

func main() {
	re = regexp.MustCompile(`(?m)(0-9a-fA-F3){1,2}`)
	re = regexp.MustCompile(`(?m)^#([0-9a-fA-F]{3}){1,2}$`)
	fmt.Printf("<%T> %+v\n", re.MatchString("#234123333"), re.MatchString("#234123333"))
}
