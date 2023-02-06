nostr-contacts
===
A command line tool for making a backup of your #nostr contact list.

(Also known as follows).

See [*nostr* protocol](https://github.com/nostr-protocol).

### features

- [x] Support Linux / Mac / Windows
- [x] Backup contacts as a flat .txt file
- [ ] Restore contacts from a backup

### [Download the release](https://github.com/jeremyd/nostr-contacts/releases) from github.

Unpack and run.

example:
```
nostr-contacts backup --relay wss://nostr21.com --relay wss://nostr-pub.wellorder.net --pubkey npub1xxx
``

### Contacts backup .txt file

Contacts will be saved in contacts-(pubkey)-(unix timestamp).txt in the current directory.

This file contains a list of the pubkeys of your follows (one per line).