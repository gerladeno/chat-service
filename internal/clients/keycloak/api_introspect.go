package keycloakclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

type IntrospectTokenResult struct {
	Exp    int      `json:"exp"`
	Iat    int      `json:"iat"`
	Aud    []string `json:"aud"`
	Active bool     `json:"active"`
}

func (t *IntrospectTokenResult) UnmarshalJSON(data []byte) error {
	var tmp struct {
		Exp    int  `json:"exp"`
		Iat    int  `json:"iat"`
		Aud    any  `json:"aud"`
		Active bool `json:"active"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return fmt.Errorf("custom UnmarshalJSON error: %w", err)
	}
	t.Exp = tmp.Exp
	t.Iat = tmp.Iat
	t.Active = tmp.Active
	if tmp.Aud == nil {
		return nil
	}
	switch typ := tmp.Aud.(type) {
	case []interface{}:
		for _, elem := range typ {
			audStr, _ := elem.(string)
			t.Aud = append(t.Aud, audStr)
		}
	case string:
		t.Aud = []string{typ}
	default:
		return errors.New("unexpected unmarshalled Aud type")
	}
	return nil
}

// IntrospectToken implements
// https://www.keycloak.org/docs/latest/authorization_services/index.html#obtaining-information-about-an-rpt
func (c *Client) IntrospectToken(ctx context.Context, token string) (*IntrospectTokenResult, error) {
	url := fmt.Sprintf("realms/%s/protocol/openid-connect/token/introspect", c.realm)
	var result IntrospectTokenResult
	resp, err := c.auth(ctx).SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(map[string]string{
			"token_type_hint": "requesting_party_token",
			"token":           token,
		}).
		SetResult(&result).
		Post(url)
	if err != nil {
		return nil, fmt.Errorf("send request to keycloak: %v", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("errored keycloak response: %v", resp.Status())
	}
	return &result, nil
}

func (c *Client) auth(ctx context.Context) *resty.Request {
	return c.cli.SetBasicAuth(c.clientID, c.clientSecret).R().SetContext(ctx)
}
