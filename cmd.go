package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"go.uber.org/dig"
)

func GetRootCommand(c *dig.Container) *cobra.Command {
	var root = &cobra.Command{
		Use:   "bench",
	}

	root.AddCommand(GetCleanCommand(c),GetRunCommand(c))
	return root
}


func GetRunCommand(c *dig.Container) *cobra.Command {
	var root = &cobra.Command{
		Use:   "run",
		Run: func(cmd *cobra.Command, args []string) {
			e := c.Invoke(func(b *Bench) {
				b.Run()
			})
			if e != nil {
				fmt.Println(e)
			}
		},
	}

	return root
}


func GetCleanCommand(c *dig.Container) *cobra.Command {
	var root = &cobra.Command{
		Use:   "clean",
		Run: func(cmd *cobra.Command, args []string) {
			e := c.Invoke(func(b *Bench) {
				b.Delete()
			})
			if e != nil {
				fmt.Println(e)
			}
		},
	}

	return root
}
