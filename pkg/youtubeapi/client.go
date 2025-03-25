package youtubeapi

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/charmbracelet/log"
)

const APIBase = "https://youtube.googleapis.com/youtube/v3"

type Client struct {
	key string
}

func New(apiKey string) *Client {
	return &Client{
		key: apiKey,
	}
}

type requestParams struct {
	Endpoint    string
	Method      string
	QueryParams url.Values
	Body        []byte
}

func (c *Client) doGetRequest(ctx context.Context, endpoint string, q url.Values) (*http.Response, error) {

	// create api request
	req, err := c.newAPIRequest(ctx, http.MethodGet, endpoint, q, nil)
	if err != nil {
		err = fmt.Errorf("error building request to %s: %w", req.URL, err)
		return nil, err
	}

	// send the request
	var resp *http.Response
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		err = fmt.Errorf("error sending request to %s: %w", req.URL, err)
		return nil, err
	}

	// error if request failed
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("recieved non-200 status: %s", resp.Status)
		return nil, err
	}

	return resp, nil

}

func (c *Client) newAPIRequest(ctx context.Context, method, endpoint string, query url.Values, body []byte) (*http.Request, error) {

	// set api key
	query.Set("key", c.key)

	// build complete url
	u := fmt.Sprintf("%s/%s?%s", APIBase, endpoint, query.Encode())

	// send request
	log.Debug("new request", "url", u)
	return http.NewRequestWithContext(ctx, method, u, bytes.NewReader(body))
}

func (c *Client) parseBody(r *http.Response) (data []byte, err error) {

	defer r.Body.Close()
	data, err = io.ReadAll(r.Body)

	if err != nil {
		log.Error("error parsing response body", "error", err, "url", r.Request.URL)
		return
	}

	// log.Debugf("response body from %s:\n%s", r.Request.URL, string(data))

	return
}
