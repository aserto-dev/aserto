package user

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/jsonx"
)

type InfoCmd struct{}

func (cmd *InfoCmd) Run(c *cc.CommonCtx) error {
	info, err := getInfo(c)
	if err != nil {
		return err
	}
	return jsonx.OutputJSON(c.UI.Output(), info)
}

type Info struct {
	Sub           string `json:"sub"`
	NickName      string `json:"nickname"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	UpdatedAt     string `json:"updated_at"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
}

// getInfo, retrieve user profile information.
func getInfo(c *cc.CommonCtx) (*Info, error) {
	req, err := http.NewRequestWithContext(c.Context, "GET", c.Auth.UserInfoURL, http.NoBody)
	if err != nil {
		return nil, err
	}

	token, err := c.Token()
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", token.Type+" "+token.Access)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var info Info

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &info)
	if err != nil {
		return nil, err
	}

	return &info, nil
}
