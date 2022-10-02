package main

import (
	"bufio"
	"fmt"
	"os"
)

func reader(o string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(o)

	i, _ := reader.ReadString('\n')
	return i
}
