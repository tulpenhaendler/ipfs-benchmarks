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

				fmt.Println("\n\n-----------------------------  CLEANUP")


				b.Delete()

				fmt.Println("\n\n-----------------------------  PROVISION")

				b.Run()

				fmt.Println("\n\n-----------------------------  BENCHMARK TIME")

				b.DoBench(100)


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
