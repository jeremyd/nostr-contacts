package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "nostr-contacts",
	Short: "nostr-contacts",
	Long:  `nostr-contacts`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("USAGE: nostr-contacts [command] [options]")
		fmt.Println("commands: backup, restore")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringSlice("relay", []string{"wss://nostr-pub.wellorder.net"}, "--relay wss://nostr-pub.wellorder.net")
	rootCmd.PersistentFlags().String("pubkey", "", "--pubkey npub1xxxx")
	rootCmd.PersistentFlags().String("file", "", "--file <filename to backup/restore>")
	viper.BindPFlag("relay", rootCmd.PersistentFlags().Lookup("relay"))
	viper.BindPFlag("pubkey", rootCmd.PersistentFlags().Lookup("pubkey"))
	viper.BindPFlag("file", rootCmd.PersistentFlags().Lookup("file"))
}
