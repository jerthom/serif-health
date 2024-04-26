package cmd

import (
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"

	"sf-takehome/anthem"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sf-takehome",
	Short: "Search an Anthem Index for all Anthem NY PPO plan file locations",
	Long: `sf-takehome is a CLI program that searches an Anthem machine readable index file table of contents
			for all Anthem NY PPO plan file locations and writes them to the specified text file. It uses a simple heuristic
			to filter for the file locations, searching for plans that have a description of "In-Network Negotiated Rates Files"
			and whose URL contains "NY_".`,
	Run: func(cmd *cobra.Command, args []string) {
		start := time.Now()
		indexPath, _ := cmd.Flags().GetString("indexPath")
		outputPath, _ := cmd.Flags().GetString("outputPath")
		anthem.ProcessIndex(indexPath, outputPath)
		end := time.Now()
		duration := end.Sub(start)
		log.Println("Time taken: ", duration)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// init sets up the flags for our CLI
func init() {
	rootCmd.Flags().StringP("indexPath", "i", "2024-04-01_anthem_index.json.gz",
		"Path to the Index json.gz file")

	rootCmd.Flags().StringP("outputPath", "o", "output.txt",
		"File path to write results to")
}
