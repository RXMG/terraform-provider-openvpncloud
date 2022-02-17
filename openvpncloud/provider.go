package openvpncloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/kaiden-rxmg/terraform-provider-openvpncloud/client"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"base_url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"client_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("OPENVPN_CLOUD_CLIENT_ID", nil),
			},
			"client_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("OPENVPN_CLOUD_CLIENT_SECRET", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"openvpncloud_network":    resourceNetwork(),
			"openvpncloud_connector":  resourceConnector(),
			"openvpncloud_route":      resourceRoute(),
			"openvpncloud_dns_record": resourceDnsRecord(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"openvpncloud_network":        dataSourceNetwork(),
			"openvpncloud_connector":      dataSourceConnector(),
			"openvpncloud_user":           dataSourceUser(),
			"openvpncloud_user_group":     dataSourceUserGroup(),
			"openvpncloud_vpn_region":     dataSourceVpnRegion(),
			"openvpncloud_network_routes": dataSourceNetworkRoutes(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	baseUrl := d.Get("base_url").(string)
	clientId := d.Get("client_id").(string)
	clientSecret := d.Get("client_secret").(string)
	openvpnClient, err := client.NewClient(baseUrl, clientId, clientSecret)
	var diags diag.Diagnostics
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create client",
			Detail:   fmt.Sprintf("Error: %v", err),
		})
		return nil, diags
	}
	return openvpnClient, nil
}
