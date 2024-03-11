# go-hetzner-dns

This is my implementation of golang support for the [Hetzner Public DNS API](https://dns.hetzner.com/api-docs). This project was part of a self-built ACME client implementation.

## Sample usage

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

	data, _ := h.ExportZoneFile("example.com")
	os.WriteFile("zone.txt", data, 0755)
}
```

### Building

This is a Go project so download and install the latest version from <https://go.dev/>.Create a new directory `test/` and create a file `main.go` with above file content in it. We may then run the following to init and build the project:

```
go mod init test
go get
go build
./test
```

## Status

As of today, the current version of Hetzner DNS Public API is `1.1.1`. So far, the following **core features** were either implemented completely or implemented partially. A lot of error handling is still missing and so are extras such as pagination.

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
- [ ] Bulk Create Records
- [ ] Bulk Update Records

### Primary Servers

- [ ] Get All Primary Servers
- [ ] Create Primary Server
- [ ] Get Primary Server
- [ ] Update Primary Server
- [ ] Delete Primary Server