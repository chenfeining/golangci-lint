package main

import (
	"github.com/chenfeining/golangci-lint/test/testdata_etc/unused_exported/lib"
)

func main() {
	lib.PublicFunc()
}
