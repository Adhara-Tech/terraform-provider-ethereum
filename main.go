package main

import (
	"github.com/Adhara-Tech/terraform-provider-ethereum/ethereum"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: ethereum.Provider})
}
