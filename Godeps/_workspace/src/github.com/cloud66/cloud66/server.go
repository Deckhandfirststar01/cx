package cloud66

import (
	"errors"
	"strings"
	"time"
)

type ServerNotification struct {
	Value  interface{} `json:"value"`
	Reason string      `json:"reason"`
	AsOf   time.Time   `json:"asof"`
}

type Server struct {
	Uid              string                        `json:"uid"`
	VendorUid        string                        `json:"vendor_uid"`
	Name             string                        `json:"name"`
	Address          string                        `json:"address"`
	Distro           string                        `json:"distro"`
	DistroVersion    string                        `json:"distro_version"`
	DnsRecord        string                        `json:"dns_record"`
	ServerType       string                        `json:"server_type"`
	ServerGroupId    int                           `json:"server_group_id"`
	Roles            []string                      `json:"server_roles"`
	StackUid         string                        `json:"stack_uid"`
	HasAgent         bool                          `json:"has_agent"`
	Params           interface{}                   `json:"params"`
	CreatedAt        time.Time                     `json:"created_at"`
	UpdatedAt        time.Time                     `json:"updated_at"`
	Region           string                        `json:"region"`
	AvailabilityZone string                        `json:"availability_zone"`
	ExtIpV4          string                        `json:"ext_ipv4"`
	HealthCode       int                           `json:"health_state"`
	UserName         string                        `json:"user_name"`
	PersonalKey      bool                          `json:"personal_key"`
	SshPrivateKey    *string                       `json:"ssh_private_key"`
	Notifications    map[string]ServerNotification `json:"notifications"`
}

func (s Server) Health() string {
	return healthStatus[s.HealthCode]
}

func (c *Client) Servers(stackUid string) ([]Server, error) {
	req, err := c.NewRequest("GET", "/stacks/"+stackUid+"/servers.json", nil)
	if err != nil {
		return nil, err
	}
	var serversRes []Server
	return serversRes, c.DoReq(req, &serversRes)
}

func (c *Client) ServerKeyInformation(stackUid string, serverUid string) (string, error) {
	server, err := c.GetServer(stackUid, serverUid, 1)
	if err != nil {
		return "", err
	}
	if server.SshPrivateKey == nil {
		return "", errors.New("SshPrivateKey not returned by server")
	}
	return *server.SshPrivateKey, nil
}

func (c *Client) GetServer(stackUid string, serverUid string, includeSshKey int) (*Server, error) {
	params := struct {
		Value int `json:"include_private_key"`
	}{
		Value: includeSshKey,
	}
	req, err := c.NewRequest("GET", "/stacks/"+stackUid+"/servers/"+serverUid+".json", params)
	if err != nil {
		return nil, err
	}
	var serverRes *Server
	return serverRes, c.DoReq(req, &serverRes)
}

func (c *Client) ServerSettings(stackUid string, serverUid string) ([]StackSetting, error) {
	req, err := c.NewRequest("GET", "/stacks/"+stackUid+"/servers/"+serverUid+"/settings.json", nil)

	if err != nil {
		return nil, err
	}

	var settingsRes []StackSetting
	return settingsRes, c.DoReq(req, &settingsRes)
}

func (c *Client) ServerSet(stackUid string, serverUid string, key string, value string) (*AsyncResult, error) {
	key = strings.Replace(key, ".", "-", -1)
	params := struct {
		Value string `json:"value"`
	}{
		Value: value,
	}
	req, err := c.NewRequest("PUT", "/stacks/"+stackUid+"/servers/"+serverUid+"/settings/"+key+".json", params)
	if err != nil {
		return nil, err
	}
	var asyncRes *AsyncResult
	return asyncRes, c.DoReq(req, &asyncRes)
}

func (s *Server) HasRole(searchRole string) bool {
	for _, role := range s.Roles {
		if role == searchRole {
			return true
		}
	}
	return false
}
