package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Connector struct {
	Id              string `json:"id,omitempty"`
	Name            string `json:"name"`
	NetworkItemId   string `json:"networkItemId"`
	NetworkItemType string `json:"networkItemType"`
	VpnRegionId     string `json:"vpnRegionId"`
	IPv4Address     string `json:"ipV4Address"`
	IPv6Address     string `json:"ipV6Address"`
}

const (
	NetworkItemTypeHost    = "HOST"
	NetworkItemTypeNetwork = "NETWORK"
)

func (c *Client) GetConnectors() ([]Connector, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/beta/connectors", c.BaseURL), nil)
	if err != nil {
		return nil, err
	}
	body, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}
	var connectors []Connector
	err = json.Unmarshal(body, &connectors)
	if err != nil {
		return nil, err
	}
	return connectors, nil
}

func (c *Client) GetConnectorById(connectorId string) (*Connector, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/beta/connectors/%s", c.BaseURL, connectorId), nil)
	if err != nil {
		return nil, err
	}
	body, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}
	var connector Connector
	err = json.Unmarshal(body, &connector)
	if err != nil {
		return nil, err
	}
	return &connector, nil
}

func (c *Client) AddNetworkConnector(connector Connector) (*Connector, error) {
	connectorJson, err := json.Marshal(connector)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/beta/connectors?networkItemId=%s&networkItemType=%s", c.BaseURL, connector.NetworkItemId, connector.NetworkItemType), bytes.NewBuffer(connectorJson))
	if err != nil {
		return nil, err
	}
	body, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}
	var conn Connector
	err = json.Unmarshal(body, &conn)
	if err != nil {
		return nil, err
	}
	return &conn, nil
}

func (c *Client) UpdateNetworkConnector(connector Connector) (*Connector, error) {
	connectorJson, err := json.Marshal(connector)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/beta/connectors/%s?networkItemId=%s&networkItemType=%s", c.BaseURL, connector.Id, connector.NetworkItemId, connector.NetworkItemType), bytes.NewBuffer(connectorJson))
	if err != nil {
		return nil, err
	}
	body, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}
	var conn Connector
	err = json.Unmarshal(body, &conn)
	if err != nil {
		return nil, err
	}
	return &conn, nil
}

func (c *Client) DeleteNetworkConnector(connectorId string, networkId string, networkItemType string) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/beta/connectors/%s?networkItemId=%s&networkItemType=%s", c.BaseURL, connectorId, networkId, networkItemType), nil)
	if err != nil {
		return err
	}
	_, err = c.DoRequest(req)
	return err
}

func (c *Client) GetConnectorsForNetwork(networkId string) ([]Connector, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/beta/connectors", c.BaseURL), nil)
	if err != nil {
		return nil, err
	}
	body, err := c.DoRequest(req)
	if err != nil {
		return nil, err
	}
	var connectors []Connector
	err = json.Unmarshal(body, &connectors)
	if err != nil {
		return nil, err
	}
	var networkConnectors []Connector
	for _, v := range connectors {
		if v.NetworkItemId == networkId {
			networkConnectors = append(networkConnectors, v)
		}
	}
	return networkConnectors, nil
}

func (c *Client) GetNetworkConnectorByName(name string, networkId string) (*Connector, error) {
	connectors, err := c.GetConnectorsForNetwork(networkId)
	if err != nil {
		return nil, err
	}
	for _, c := range connectors {
		if c.Name == name {
			return &c, nil
		}
	}
	return nil, nil
}
