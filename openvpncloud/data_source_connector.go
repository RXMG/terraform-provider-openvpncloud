package openvpncloud

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/kaiden-rxmg/terraform-provider-openvpncloud/client"
)

func dataSourceConnector() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceConnectorRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"network_item_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"network_item_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpn_region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_v4_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_v6_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceConnectorRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	var diags diag.Diagnostics
	connectorName := d.Get("name").(string)
	networkitemId := d.Get("network_item_id").(string)
	connector, err := c.GetNetworkConnectorByName(connectorName, networkitemId)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	if connector == nil {
		return append(diags, diag.Errorf("Connector with name %s was not found", connectorName)...)
	}
	d.Set("name", connector.Name)
	d.Set("network_item_id", connector.NetworkItemId)
	d.Set("network_item_type", connector.NetworkItemType)
	d.Set("vpn_region_id", connector.VpnRegionId)
	d.Set("ip_v4_address", connector.IPv4Address)
	d.Set("ip_v6_address", connector.IPv6Address)
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return diags
}
