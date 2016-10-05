package marathon

import "net/url"

func (c *MarathonClient) GetMarathonInfo() (*MarathonInfo, error) {
	info := new(MarathonInfo)

	resp := c.http.HttpGet(c.marathonUrl(API_INFO), info)
	if resp.Error != nil {
		return nil, resp.Error
	}
	return info, nil
}

func (c *MarathonClient) GetCurrentLeader() (*LeaderInfo, error) {
	info := new(LeaderInfo)

	resp := c.http.HttpGet(c.marathonUrl(API_LEADER), info)
	if resp.Error != nil {
		return nil, resp.Error
	}
	return info, nil
}

func (c *MarathonClient) AbdicateLeader() (*Message, error) {
	msg := new(Message)
	resp := c.http.HttpDelete(c.marathonUrl(API_LEADER), nil, msg)
	if resp.Error != nil {
		return nil, resp.Error
	}
	return msg, nil
}

func (c *MarathonClient) Ping() (*MarathonPing, error) {
	resp := c.http.HttpGet(c.marathonUrl(API_PING), nil)
	if resp.Error != nil {
		return nil, resp.Error
	}
	host := c.getHost()
	if u, err := url.Parse(host); err == nil {
		host = u.Host
	}

	return &MarathonPing{Host: host, Elapsed: resp.Elapsed}, nil
}
