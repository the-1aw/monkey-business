package cmd

import (
	"fmt"
	"os"
	"os/user"

	"github.com/spf13/cobra"
	"github.com/the-1aw/monkey-business/repl"
)

var walkCmd = &cobra.Command{
	Use:   "walk",
	Short: "Run monkey code with tree walking interpreter",
	Long: `
This command runs a monkey tree walking interpreter REPL (Read-Eval-Print-Loop).
Walk command is ~3x slower than the run command
`,
	Run: func(cmd *cobra.Command, args []string) {
		user, err := user.Current()
		if err != nil {
			panic(err)
		}
		fmt.Printf("Hello %s! This is the Monkey programming language!\n", user.Username)
		fmt.Printf("Feel free to type in commands\n")
		repl.StartInterpreter(os.Stdin, os.Stdout)
	},
}

func init() {
	rootCmd.AddCommand(walkCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// walkCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// walkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
