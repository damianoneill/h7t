package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func generateMarkdown() {
	err := doc.GenMarkdownTree(rootCmd, "./docs/")
	if err != nil {
		log.Fatal(err)
	}
}

// docsCmd represents the docs command
var docsCmd = &cobra.Command{
	Use:   "docs",
	Short: "Generate Markdown for the commands in " + rootCmd.Use,
	Long:  `For ` + rootCmd.Use + ` generate Markdown Documents for each of the commands and write them to a folder named ./docs`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Writing command descriptions to ./docs")
		generateMarkdown()
	},
}

func init() {
	rootCmd.AddCommand(docsCmd)
}
