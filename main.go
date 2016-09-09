package main

import (
	"fmt"
	"os"

	"github.com/yroffin/jarvis-go-ext/cmd"
)

// main function
func main() {

	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

}
