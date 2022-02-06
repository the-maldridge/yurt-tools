package main

import (
	"github.com/the-maldridge/yurt-tools/internal/cmdlets"
	// additional custom filters
	_ "github.com/the-maldridge/yurt-tools/internal/http/filters"
)

func main() {
	cmdlets.Entrypoint()
}
