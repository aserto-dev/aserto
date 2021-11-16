package user

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/aserto-dev/aserto/pkg/auth0"
	"github.com/aserto-dev/aserto/pkg/cc"
	"github.com/aserto-dev/aserto/pkg/jsonx"
)

type InfoCmd struct{}

func (cmd *InfoCmd) Run(c *cc.CommonCtx) error {
	info, err := getInfo(c)
	if err != nil {
		return err
	}
	return jsonx.OutputJSON(c.OutWriter, info)
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
	env := c.Environment()
	req, err := http.NewRequestWithContext(c.Context, "GET", auth0.GetSettings(env).UserInfoURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", c.Token().Type+" "+c.Token().Access)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var info Info

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &info)
	if err != nil {
		return nil, err
	}

	return &info, nil
}
