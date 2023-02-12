package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backup contact list",
	Long:  `Backup contact list`,
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
		} else if isSetPub {
			pub = decodePub(pk)
		} else if cmdOptPub {
			pub = decodePub(viper.GetString("pubkey"))
		}

		ctx := context.Background()

		relays := viper.GetStringSlice("relay")
		if !viper.IsSet("relay") {
			fmt.Printf("--relays not specified, using defaults: %s\n", relays)
		}

		filters := []nostr.Filter{
			{
				Kinds:   []int{3},
				Limit:   1,
				Authors: []string{pub},
			},
		}

		var allSubs []*nostr.Subscription
		var allRelays []*nostr.Relay
		var allContact []string

		contactsChannel := make(chan string)
		relaysChannel := make(chan *nostr.Relay)
		subsChannel := make(chan *nostr.Subscription)

		fmt.Println("Waiting up to 10 seconds for all relays to respond...")

		go func() {
			for {
				r := <-relaysChannel
				allRelays = append(allRelays, r)
				relayURL := r.URL
				log(fmt.Sprintf("> %-30s connected.", relayURL))
			}
		}()

		go func() {
			for {
				s := <-subsChannel
				allSubs = append(allSubs, s)
				relayURL := s.Relay.URL
				<-s.EndOfStoredEvents
				log(fmt.Sprintf("> %-30s status: query complete.", relayURL))
			}
		}()

		for _, r := range relays {
			relay, err := nostr.RelayConnect(ctx, r)
			if err != nil {
				log(fmt.Sprintf("failed initial connection to relay: %s, %s; skipping relay.", r, err))
				continue
			}
			sub := relay.Subscribe(ctx, filters)

			go func() {
				subsChannel <- sub
				relaysChannel <- relay
			}()

			go func() {
				for n := range relay.Notices {
					log(fmt.Sprintf("notice from %s: %s", relay.URL, n))
				}
			}()

			go func() {
				for ev := range sub.Events {
					if ev.Kind == 3 {
						// Contact List
						pTags := []string{"p"}
						allPTags := ev.Tags.GetAll(pTags)
						log(fmt.Sprintf("(%d) contacts found on relay: %-30s", len(allPTags), relay.URL))
						for _, tag := range allPTags {
							contactsChannel <- tag[1]
						}
					}
				}
			}()
		}

		unixTime := time.Now().Unix()

		fileName := fmt.Sprintf("contacts-%s-%d.txt", pub, unixTime)

		go func() {
			for {
				newContact := <-contactsChannel
				allContact = append(allContact, newContact)
			}
		}()

		startTime := time.Now()

		for {
			time.Sleep(5 * time.Second)
			timeout := time.Now().After(startTime.Add(20 * time.Second))
			if timeout {
				if timeout {
					log(fmt.Sprintln("timeout (20s) reached."))
				}

				log(fmt.Sprintf("closing connections."))
				for _, sub := range allSubs {
					log(fmt.Sprintf("closing subscription: %v on %s", sub.Filters, sub.Relay.URL))
					sub.Unsub()
				}
				for _, relay := range allRelays {
					log(fmt.Sprintf("closing relay: %s ", relay.URL))
				}

				allContact = removeDupes(allContact)

				log(fmt.Sprintf("...found %d contacts", len(allContact)))

				f, err := os.Create(fileName)
				if err != nil {
					log(fmt.Sprintf("Error opening file %s, %s; exiting", err, fileName))
					os.Exit(1)
				}
				defer f.Close()
				for _, c := range allContact {
					f.WriteString(c + "\n")
				}
				break
			}
		}

		fmt.Printf("Done. %d contacts saved in %s ", len(allContact), fileName)
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)
}
