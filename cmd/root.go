package cmd

import (
    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "myapp",
    Short: "A task management application",
}

func Execute() error {
    return rootCmd.Execute()
}

func init() {
    rootCmd.AddCommand(CreateCmd)
    rootCmd.AddCommand(ShowCmd)
    rootCmd.AddCommand(UpdateCmd)
    rootCmd.AddCommand(DeleteCmd)
}
