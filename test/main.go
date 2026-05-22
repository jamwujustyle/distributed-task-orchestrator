package main

import (
	"fmt"
	"strings"

	"github.com/jamwujustyle/logger"
)

func main() {
	logger.InitLogger(1 == 0)

	str := "hello what the fuck"

	fmt.Printf("%#v", strings.Fields(str))
}
