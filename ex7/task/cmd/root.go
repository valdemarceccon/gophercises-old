package cmd

import "github.com/spf13/cobra"

// RootCmd command for the task app
var RootCmd = &cobra.Command{
	Use:   "task",
	Short: "Task is a cli task mananger",
}
