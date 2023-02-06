package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
)

func removeDupes(stringSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func main() {
	sk, isSet := os.LookupEnv("NOSTR_PRIVATE")
	pk, isSetPub := os.LookupEnv("NOSTR_PUBLIC")

	if !isSet && !isSetPub {
		fmt.Println("please set the environment variable for NOSTR_PRIVATE or NOSTR_PUBLIC to gather contacts.")
		os.Exit(1)
	}
	var pub string
	var npub string
	var privatek string

	if isSet {
		if sk[:4] == "nsec" {
			_, v, _ := nip19.Decode(sk)
			privatek = v.(string)
		} else {
			privatek = sk
		}
		pub, _ := nostr.GetPublicKey(privatek)
		npub, _ = nip19.EncodePublicKey(pub)
	} else if isSetPub {
		if pk[:4] == "npub" {
			_, v, _ := nip19.Decode(pk)
			pub = v.(string)
			npub = pk
		} else {
			pub = pk
			npub, _ = nip19.EncodePublicKey(pub)
		}
	}

	fmt.Println("pubkey:", pub)
	fmt.Println("npub: ", npub)

	ctx := context.Background()

	relays := []string{
		"wss://nostr-dev.wellorder.net",
		"wss://nostr-pub.wellorder.net",
		"wss://nostr21.com",
		"wss://no.str.cr",
		"wss://nostr.cercatrova.me",
		"wss://nostr.oxtr.dev",
		"wss://nostr.p2sh.co",
		"wss://nostr.pleb.network",
		"wss://relay.orangepill.dev",
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

	fmt.Println("Waiting up to 10 seconds for all relays to respond...")

	for _, r := range relays {
		relay, err := nostr.RelayConnect(ctx, r)
		if err != nil {
			fmt.Printf("failed initial connection to relay: %s, %s; skipping relay\n", r, err)
			continue
		}
		sub := relay.Subscribe(ctx, filters)
		allSubs = append(allSubs, sub)
		allRelays = append(allRelays, relay)
		go func() {
			for ev := range sub.Events {
				if ev.Kind == 3 {
					// Contact List
					pTags := []string{"p"}
					allPTags := ev.Tags.GetAll(pTags)
					for _, tag := range allPTags {
						contactsChannel <- tag[1]
					}
				}
			}
		}()
	}

	eoseCompleted := false
	go func() {
		for i, sub := range allSubs {
			relayURL := sub.Relay.URL
			<-sub.EndOfStoredEvents
			fmt.Printf("Contact list received from %-30s (%d of %d).\n", relayURL, i+1, len(allSubs))
		}
		eoseCompleted = true

	}()

	unixTime := time.Now().Unix()

	fileName := fmt.Sprintf("contacts-%d.txt", unixTime)

	go func() {
		for {
			newContact := <-contactsChannel
			allContact = append(allContact, newContact)
		}
	}()

	startTime := time.Now()

	for {
		time.Sleep(5 * time.Second)
		timeout := time.Now().After(startTime.Add(10 * time.Second))
		if eoseCompleted || timeout {
			if timeout {
				fmt.Println("timeout (10s) reached.")
			}
			fmt.Println("closing connections.")
			for _, sub := range allSubs {
				sub.Unsub()
			}
			for _, relay := range allRelays {
				relay.Close()
			}

			allContact = removeDupes(allContact)

			fmt.Printf("...found %d contacts\n", len(allContact))

			f, err := os.Create(fileName)
			if err != nil {
				fmt.Printf("Error opening file %s, %s; exiting", err, fileName)
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
}
