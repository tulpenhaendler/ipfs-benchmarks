package main

import (
	"fmt"
	"go.uber.org/dig"
	"os"
)

func main(){
	c := dig.New()
	c.Provide(NewConfig)

	rootCmd := GetRootCommand(c)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}