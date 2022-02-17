package main

import (
	"github.com/kaiden-rxmg/terraform-provider-openvpncloud/openvpncloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return openvpncloud.Provider()
		},
	})
}
