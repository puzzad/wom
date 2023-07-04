package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func NewImportCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "import",
		Short: "Imports adventures",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Importing")
		},
	}
}
