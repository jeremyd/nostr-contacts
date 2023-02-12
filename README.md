nostr-contacts
===
A command line tool for making a backup of your #nostr contact list.

(Also known as follows).

See [*nostr* protocol](https://github.com/nostr-protocol).

### features

- [x] Support Linux / Mac / Windows
- [x] Backup contacts as a flat .txt file
- [x] Restore contacts from a backup

### [Download the release](https://github.com/jeremyd/nostr-contacts/releases) from github.

Unpack and run.

### help:
```
nostr-contacts --help
```

### example backing up contacts:
```
nostr-contacts backup --relay wss://nostr21.com --relay wss://nostr-pub.wellorder.net --pubkey npub1xxx
```

### example restoring contacts:
```
export NOSTR_PRIVATE=<your private key>
# this will prompt for confirmation before broadcasting
nostr-contacts restore --file <path to backup.txt> --relay wss://nostr-pub.wellorder.net
```

### Contacts Backup .txt file

Contacts will be saved in contacts-(pubkey)-(unix timestamp).txt in the current directory.

This file contains a list of the pubkeys of your follows (one per line).

### Contacts Restore
Contacts can be restored from any backup by specifying --file on the command line and the list of relays that you'd like to broadcast to.