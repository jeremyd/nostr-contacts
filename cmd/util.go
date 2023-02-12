package cmd

import (
	"fmt"
	"os"

	"github.com/nbd-wtf/go-nostr/nip19"
)

//var logOut = bufio.NewWriter(os.Stdout)

func log(message string) {
	fmt.Println(message)
	//logOut.WriteString(fmt.Sprintln(message))
	//logOut.Flush()
}

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

func decodePub(pk string) string {
	var pub string
	var npub string
	if pk[:4] == "npub" {
		_, v, err := nip19.Decode(pk)
		if err != nil {
			log(fmt.Sprintf("could not decode pubkey for %s", pk))
			os.Exit(1)
		}
		pub = v.(string)
		npub = pk
	} else {
		pub = pk
		npub, _ = nip19.EncodePublicKey(pub)
	}
	fmt.Println("pubkey:", pub)
	fmt.Println("npub:", npub)
	return pub
}
