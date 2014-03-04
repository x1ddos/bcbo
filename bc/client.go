package bc

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"path"
)

const (
	paramApiKey      = "apikey"
	paramSearchPath  = "path"
	paramSearchValue = "value"
)

// client is the default implementation of Client.
type client struct {
	apiroot string
}

// NewClient creates a new instance of Beancounter client.
// apiroot should point to the root of Beancounter Platform API interface.
func NewClient(apiroot string) (Client, error) {
	u, err := url.Parse(apiroot)
	if err != nil {
		return nil, err
	}

	u.RawQuery = ""
	u.Fragment = ""
	apiroot = u.String()
	if apiroot[len(apiroot)-1] != '/' {
		apiroot += "/"
	}

	return &client{apiroot}, nil
}

func (c *client) ListUsers(cred *Cred) ([]*User, error) {
	bcUsersList := &bcApiUsersList{}
	if err := c.doGetRequest(cred, "user/all", nil, bcUsersList); err != nil {
		return nil, err
	}
	if bcUsersList.Status != "OK" {
		return nil, errors.New(bcUsersList.Status + ": " + bcUsersList.Message)
	}

	list := make([]*User, 0, len(bcUsersList.Object))
	for _, bcUser := range bcUsersList.Object {
		list = append(list, toUserResource(bcUser))
	}

	return list, nil
}

func (c *client) GetUserProfile(cred *Cred, userId string) (*UserProfile, error) {
	bcProfile := &bcApiUserProfile{}
	resource := path.Join("user", userId, "profile")
	if err := c.doGetRequest(cred, resource, nil, bcProfile); err != nil {
		return nil, err
	}
	if bcProfile.Status != "OK" {
		return nil, errors.New(bcProfile.Status + ": " + bcProfile.Message)
	}

	return toUserProfileResource(bcProfile), nil
}

func (c *client) ListActivities(cred *Cred, userId string, page int) ([]*Activity, error) {
	query := url.Values{}
	query.Set(paramSearchPath, "user.username")
	query.Set(paramSearchValue, userId)
	bcActivities := &bcApiActivitiesList{}
	if err := c.doGetRequest(cred, "activities/search", query, bcActivities); err != nil {
		return nil, err
	}
	if bcActivities.Status != "OK" {
		return nil, errors.New(bcActivities.Status + ": " + bcActivities.Message)
	}

	list := make([]*Activity, 0, len(bcActivities.Object))
	for _, item := range bcActivities.Object {
		list = append(list, toActivityResource(&item.ApiActivity))
	}

	return list, nil
}

func (c *client) formatUrl(cred *Cred, resource string, query url.Values) string {
	if query == nil {
		query = url.Values{}
	}
	if cred != nil && cred.ApiKey != "" {
		query.Set(paramApiKey, cred.ApiKey)
	}
	if resource[0] == '/' {
		resource = resource[1:]
	}
	return c.apiroot + resource + "?" + query.Encode()
}

func (c *client) doGetRequest(cred *Cred, resource string, query url.Values, dst interface{}) error {
	resourceUrl := c.formatUrl(cred, resource, query)
	resp, err := http.Get(resourceUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(dst)
}
