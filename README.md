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

set your PUBLIC or PRIVATE key as an environment variable.
```
# public key supports npub or hex
export NOSTR_PUBLIC=npub1xxxxx

# private key supports nsec or hex
export NOSTR_PRIVATE=nsec1xxxxx
```

Unpack and run.

### Contacts backup .txt file

Contacts will be saved in contacts-(unix timestamp).txt

This file contains a list of the pubkeys of your follows (one per line).