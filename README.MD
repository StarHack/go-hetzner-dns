# go-hetzner-dns

This is my implementation of golang support for the Hetzner DNS Public API. This project was part of a self-built ACME client implementation.

## Sample usage

```
package main

import "fmt"
import "github.com/StarHack/go-hetzner-dns"

func main() {
	h := Hetzner{APIKey: "<your_api_key>"}
	zones, err := h.FindAllZones()

	if err == nil {
		panic(err)
	}

	for _, zone := range zones {
		fmt.Println(zone.Name)
	}

}

```

## Status

As of today, the current version of Hetzner DNS Public API is `1.1.1`. So far, the following features were implemented or implemented partially.

### Authentication

- [x] Auth-API-Token

### Zones

- [x] Get All Zones
- [x] Create Zone
- [x] Get Zone
- [x] Update Zone
- [x] Delete Zone
- [ ] Import Zone file plain
- [ ] Export Zone file
- [ ] Validate Zone file plain

### Records

- [x] Get All Records
- [x] Create Record
- [x] Get Record
- [x] Update Record
- [x] Delete Record
- [ ] Bulk Create Records
- [ ] Bulk Update Records

### Primary Servers

- [ ] Get All Primary Servers
- [ ] Create Primary Server
- [ ] Get Primary Server
- [ ] Update Primary Server
- [ ] Delete Primary Server