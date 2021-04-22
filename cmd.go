package main

import (
	"github.com/spf13/cobra"
	"go.uber.org/dig"
)

func GetRootCommand(c *dig.Container) *cobra.Command {
	var root = &cobra.Command{
		Use:   "ipfs-bench",
	}

	root.AddCommand(GetCleanCommand(c),GetRunCommand(c))
	return root
}


func GetRunCommand(c *dig.Container) *cobra.Command {
	var root = &cobra.Command{
		Use:   "run",
	}

	return root
}


func GetCleanCommand(c *dig.Container) *cobra.Command {
	var root = &cobra.Command{
		Use:   "clean",
	}

	return root
}
