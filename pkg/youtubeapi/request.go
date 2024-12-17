package youtubeapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/charmbracelet/log"
	"gitlab.com/andyinabox/yesterdays-news-downloader/pkg/youtubeapi/response"
)

type requestParams struct {
	Endpoint    string
	Method      string
	QueryParams url.Values
	Body        []byte
}

func (c *Client) doGetRequest(ctx context.Context, endpoint string, q url.Values) (*response.Success, error) {
	req, err := c.newRequest(ctx, requestParams{
		Endpoint:    endpoint,
		Method:      http.MethodGet,
		QueryParams: q,
	})
	if err != nil {
		return nil, err
	}

	var resp *http.Response
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("recieved non-200 status: %s", resp.Status)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	log.Debug(string(body))

	data := response.Success{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (c *Client) newRequest(ctx context.Context, params requestParams) (*http.Request, error) {
	params.QueryParams.Add("key", c.key)
	u := fmt.Sprintf("%s/%s?%s", APIBase, params.Endpoint, params.QueryParams.Encode())
	log.Debug("new request", "url", u)
	return http.NewRequestWithContext(ctx, params.Method, u, bytes.NewReader(params.Body))
}
