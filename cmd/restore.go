package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore contact list",
	Long:  `Restore contact list`,
	Run: func(cmd *cobra.Command, args []string) {

		sk, isSet := os.LookupEnv("NOSTR_PRIVATE")
		pk, isSetPub := os.LookupEnv("NOSTR_PUBLIC")

		cmdOptPub := viper.IsSet("pubkey")

		if !isSet && !isSetPub && !cmdOptPub {
			fmt.Println("please set the environment variable for NOSTR_PRIVATE or NOSTR_PUBLIC or --pubkey")
			os.Exit(1)
		}
		var pub string
		var privatek string

		if isSet {
			if sk[:4] == "nsec" {
				_, v, _ := nip19.Decode(sk)
				privatek = v.(string)
			} else {
				privatek = sk
			}
			pub, _ = nostr.GetPublicKey(privatek)
			decodePub(pub)
		} else if isSetPub {
			pub = decodePub(pk)
		} else if cmdOptPub {
			pub = decodePub(viper.GetString("pubkey"))
		} else {
			fmt.Println("please set the environment variable for NOSTR_PRIVATE (private key required for restore)")
			os.Exit(1)
		}

		ctx := context.Background()
		currentTime := time.Now()

		relays := viper.GetStringSlice("relay")
		if !viper.IsSet("relay") {
			fmt.Printf("--relays not specified, using defaults: %s\n", relays)
		}

		var allRelays []*nostr.Relay

		relaysChannel := make(chan *nostr.Relay)

		// load keys from file
		// create nostr.Tags from keys
		var followTags []nostr.Tag

		// open file for reading
		if !viper.IsSet("file") {
			fmt.Println("please specify a file to restore from with --file")
		}
		restoreFile := viper.GetString("file")
		file, err := os.Open(restoreFile)
		if err != nil {
			log(fmt.Sprintf("failed to open file: %s, %s; aborting.", restoreFile, err))
			os.Exit(1)
		}
		defer file.Close()

		// create a scanner and read the file line by line
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			tag := nostr.Tag{"p", line}
			followTags = append(followTags, tag)
		}

		// create a new event
		ev := nostr.Event{
			PubKey:    pub,
			CreatedAt: currentTime,
			Kind:      nostr.KindContactList,
			Tags:      followTags,
			Content:   "",
		}

		// calling Sign sets the event ID field and the event Sig field
		ev.Sign(privatek)

		fmt.Printf("Are you SURE you want to broadcast your contacts list of (%d)follows to (%d)relays? (y/n)\n", len(ev.Tags), len(relays))
		var ok []byte = make([]byte, 1)
		os.Stdin.Read(ok)
		for _, c := range ok {
			if c != 'y' {
				fmt.Println("Exiting..")
				os.Exit(1)
			}
		}

		//log(fmt.Sprintf("%v", ev))

		go func() {
			for {
				r := <-relaysChannel
				allRelays = append(allRelays, r)
				//relayURL := r.URL
				//log(fmt.Sprintf("> %-30s connected.", relayURL))
			}
		}()

		for _, r := range relays {
			relay, err := nostr.RelayConnect(ctx, r)
			if err != nil {
				log(fmt.Sprintf("failed initial connection to relay: %s, %s; skipping relay.", r, err))
				continue
			}

			go func() {
				relaysChannel <- relay
			}()

			go func() {
				for n := range relay.Notices {
					log(fmt.Sprintf("NOTICE from %s: %s", relay.URL, n))
				}
			}()

			// broadcast
			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			result := relay.Publish(ctx, ev)
			log(fmt.Sprintf("broadcast follow list (%d follows) to %s: result -> %s", len(ev.Tags), relay.URL, result))
		}

		/*
			log("closing connections.")
			for _, relay := range allRelays {
				log(fmt.Sprintf("closing relay: %s ", relay.URL))
			}
		*/

		log("Done broadcasting.")
	},
}

func init() {
	rootCmd.AddCommand(restoreCmd)
}
