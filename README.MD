# go-hetzner-dns

This is my implementation of golang support for the [Hetzner Public DNS API](https://dns.hetzner.com/api-docs). This project was part of a self-built ACME client implementation and it was needed to update DNS records of domains during the ACME Challenge.

Visit [this page](https://dns.hetzner.com/settings/api-token) to generate a Hetzner Public DNS API Token.

## Sample usage

Refer to the [wiki tab](https://github.com/StarHack/go-hetzner-dns/wiki) in this repository for further examples.

```go
package main

import (
  "fmt"
  "os"

  "github.com/StarHack/go-hetzner-dns"
)

func main() {
  h := hetzner.Hetzner{APIKey: "<your_api_key>"}
  zones, err := h.FindAllZones()

  if err != nil {
    panic(err)
  }

  for _, zone := range zones {
    fmt.Println(zone.Name)
  }

  zoneId, err := h.FindZoneID("example.com")

  if err != nil {
    panic(err)
  }

  h.CreateOrUpdateRecord(zoneId, "A", "www", "1.2.3.4")
  h.CreateRecord(zoneId, "TXT", "_acme-challenge", "ABCD")
  h.UpdateRecord(...)
  h.DeleteRecord(...)

  data, _ := h.ExportZoneFile("example.com")
  os.WriteFile("zone.txt", data, 0755)
}
```

### Building

This is a Go project so download and install the latest version from <https://go.dev/>. Create a new directory `test/` and create a file `main.go` with above file content in it. We may then run the following to init and build the project:

```
go mod init test
go get
go build
./test
```

## Status

As of today, the current version of Hetzner DNS Public API is `1.1.1`. So far, the following **core features** were implemented. A couple of features such as pagination are still missing but the basic functionality for all of the provided endpoints is there.

### Authentication

- [x] Auth-API-Token

### Zones

- [x] Get All Zones
- [x] Create Zone
- [x] Get Zone
- [x] Update Zone
- [x] Delete Zone
- [x] Import Zone file plain
- [x] Export Zone file
- [x] Validate Zone file plain

### Records

- [x] Get All Records
- [x] Create Record
- [x] Get Record
- [x] Update Record
- [x] Delete Record
- [x] Bulk Create Records
- [x] Bulk Update Records

### Primary Servers

**The following endpoints were implemented but not tested yet.**

- [x] Get All Primary Servers
- [x] Create Primary Server
- [x] Get Primary Server
- [x] Update Primary Server
- [x] Delete Primary Server
