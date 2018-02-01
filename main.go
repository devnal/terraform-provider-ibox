package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/terraform-providers/terraform-provider-ibox/ibox"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: ibox.Provider})
}
