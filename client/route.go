package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type Route struct {
	Id          string `json:"id"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Subnet      string `json:"subnet"`
	Domain      string `json:"domain"`
	Value       string `json:"value"`
	NetworkId   string
}

const (
	RouteTypeIPV4   = "IP_V4"
	RouteTypeIPV6   = "IP_V6"
	RouteTypeDomain = "DOMAIN"
)

func (c *Client) CreateRoute(route Route, networkId string) (*Route, error) {
	routeJson, err := json.Marshal(route)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/beta/networks/%s/routes", c.BaseURL, networkId), bytes.NewBuffer(routeJson))
	if err != nil {
		return nil, err
	}
	body, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}
	var r Route
	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (c *Client) GetRoutes(networkId string) ([]Route, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/beta/networks/%s/routes", c.BaseURL, networkId), nil)
	if err != nil {
		return nil, err
	}
	body, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}
	var routes []Route
	err = json.Unmarshal(body, &routes)
	if err != nil {
		return nil, err
	}
	return routes, nil
}

func (c *Client) UpdateRoute(route Route, networkId string) (*Route, error) {
	routeJson, err := json.Marshal(route)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/api/beta/networks/%s/routes/%s", c.BaseURL, networkId, route.Id), bytes.NewBuffer(routeJson))
	if err != nil {
		return nil, err
	}
	body, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}
	var r Route
	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func (c *Client) DeleteRoute(networkId string, routeId string) error {
	return c.DeleteRouteHelper(networkId, routeId, 3)
}

func (c *Client) DeleteRouteHelper(networkId string, routeId string, retries int) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/beta/networks/%s/routes/%s", c.BaseURL, networkId, routeId), nil)
	if err != nil {
		return err
	}

	_, err = c.DoRequest(req)
	if err != nil {
		return err
	}

	r, err := c.GetNetworkRoute(networkId, routeId)
	if r != nil {
		if retries > 0 {
			return c.DeleteRouteHelper(networkId, routeId, retries-1)
		}
		return errors.New("could not delete route")
	}
	return err
}

func (c *Client) GetNetworkRoute(networkId string, routeId string) (*Route, error) {
	routes, err := c.GetRoutes(networkId)
	if err != nil {
		return nil, err
	}
	for _, r := range routes {
		if r.Id == routeId {
			return &r, nil
		}
	}
	return nil, nil
}

func (c *Client) GetRouteById(routeId string) (*Route, error) {
	networks, err := c.GetNetworks()
	if err != nil {
		return nil, err
	}
	for _, n := range networks {
		r, err := c.GetNetworkRoute(n.Id, routeId)
		if err != nil {
			return nil, err
		}
		if r != nil {
			r.NetworkId = n.Id
			return r, nil
		}
	}
	return nil, nil
}
