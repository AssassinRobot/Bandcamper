package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "bandcamp-downloader",
	Short: "A simple CLI tool to download music from Bandcamp",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to Bandcamp Downloader!")
	},
}

func init() {
	rootCmd.AddCommand(bandCmd)
	rootCmd.AddCommand(urlCmd)
	rootCmd.AddCommand(wishlistCmd)
	rootCmd.PersistentFlags().BoolP("force", "f", false, "Force re-download of files")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
