package openvpncloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/kaiden-rxmg/terraform-provider-openvpncloud/client"
)

func resourceRoute() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRouteCreate,
		ReadContext:   resourceRouteRead,
		UpdateContext: resourceRouteUpdate,
		DeleteContext: resourceRouteDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"network_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{client.RouteTypeIPV4, client.RouteTypeIPV6, client.RouteTypeDomain}, false),
			},
			"value": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceRouteCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	var diags diag.Diagnostics
	networkId := d.Get("network_id").(string)
	routeDescription := d.Get("description").(string)
	routeType := d.Get("type").(string)
	routeValue := d.Get("value").(string)
	r := client.Route{
		Type:        routeType,
		Description: routeDescription,
		Value:       routeValue,
	}
	route, err := c.CreateRoute(r, networkId)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	d.SetId(route.Id)
	if routeType == client.RouteTypeIPV4 || routeType == client.RouteTypeIPV6 {
		d.Set("value", route.Subnet)
	} else if routeType == client.RouteTypeDomain {
		d.Set("value", route.Domain)
	}
	return diags
}

func resourceRouteRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	var diags diag.Diagnostics
	routeId := d.Id()
	networkId := d.Get("network_id").(string)
	r, err := c.GetNetworkRoute(networkId, routeId)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	if r == nil {
		d.SetId("")
	} else {
		d.Set("type", r.Type)
		if r.Type == client.RouteTypeIPV4 || r.Type == client.RouteTypeIPV6 {
			d.Set("value", r.Subnet)
		} else if r.Type == client.RouteTypeDomain {
			d.Set("value", r.Domain)
		}
	}
	return diags
}

func resourceRouteUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	var diags diag.Diagnostics
	_, networkId := d.GetChange("network_id")
	_, description := d.GetChange("description")
	_, routeType := d.GetChange("type")
	_, routeValue := d.GetChange("value")

	r := client.Route{
		Id:          d.Id(),
		Description: description.(string),
		Type:        routeType.(string),
		Value:       routeValue.(string),
	}
	route, err := c.UpdateRoute(r, networkId.(string))
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	d.SetId(route.Id)
	if routeType == client.RouteTypeIPV4 || routeType == client.RouteTypeIPV6 {
		d.Set("value", route.Subnet)
	} else if routeType == client.RouteTypeDomain {
		d.Set("value", route.Domain)
	}
	return diags
}

func resourceRouteDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)
	var diags diag.Diagnostics
	routeId := d.Id()
	networkId := d.Get("network_id").(string)
	err := c.DeleteRoute(networkId, routeId)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	return diags
}
