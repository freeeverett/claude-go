package claude

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

func New(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		domain: DefaultDomain,
	}
}

func (c *Client) SetDomain(domain string) {
	// SetDomain 设置请求的域名，默认为 https://api.anthropic.com/
	c.domain = domain
}

func (c *Client) CreateMessage(ctx context.Context, in *RequestMessageContent) (*ResponseMessageContent, error) {
	if in == nil {
		return nil, errors.New("input is nil")
	}
	b, _ := json.Marshal(in)
	cli := http.Client{}
	req, err := http.NewRequest("POST", c.domain+"v1/messages", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	in.Stream = false

	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", anthropicVersion)
	req.Header.Set("content-type", "application/json")

	select {
	case <-ctx.Done():
		return nil, errors.New("context is done")
	default:
	}
	r, errDo := cli.Do(req)
	if errDo != nil {
		return nil, errDo
	}
	if r.StatusCode != 200 {
		var errResult = ResponseError{}
		_ = json.NewDecoder(r.Body).Decode(&errResult)
		errResult.Code = r.StatusCode
		return nil, &errResult
	}
	var res ResponseMessageContent
	return &res, json.NewDecoder(r.Body).Decode(&res)
}

func (c *Client) CreateSimpleMessage(ctx context.Context, text string) (string, error) {
	if text == "" {
		return "", errors.New("text is empty")
	}
	r, err := c.CreateMessage(ctx, &RequestMessageContent{
		Model:     DefaultModel,
		MaxTokens: DefaultMaxTokens,
		Messages: []*Message{
			{
				Role:    RoleUser,
				Content: text,
			},
		},
	})
	if err != nil {
		return "", err
	}
	if len(r.Content) > 0 && r.Content[0].Type == "text" {
		return r.Content[0].Text, nil
	}
	return "", errors.New("text is empty")
}

func (c *Client) CreateMessageStream(ctx context.Context, in *RequestMessageContent) (<-chan ResponseMessageStream, error) {
	if in == nil {
		return nil, errors.New("input is nil")
	}

	in.Stream = true

	b, _ := json.Marshal(in)
	cli := http.Client{}

	req, err := http.NewRequest("POST", c.domain+"v1/messages", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", anthropicVersion)
	req.Header.Set("content-type", "application/json")
	req.Header.Set("anthropic-beta", "messages-2023-12-15")

	r, errDo := cli.Do(req)
	if errDo != nil {
		return nil, errDo
	}
	if r.StatusCode != http.StatusOK {
		var errResult = ResponseError{}
		_ = json.NewDecoder(r.Body).Decode(&errResult)
		errResult.Code = r.StatusCode
		return nil, &errResult
	}
	reader := bufio.NewReader(r.Body)
	//
	ch := make(chan ResponseMessageStream)
	go func() {
		defer close(ch)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				line, err := reader.ReadBytes('\n')
				if err != nil {
					return
				}
				if len(line) > 6 && string(line[:5]) == "data:" {
					var res ResponseMessageStream
					_ = json.Unmarshal(line[6:], &res)
					ch <- res
				}
			}
		}
	}()
	return ch, nil
}

func (c *Client) CreateSimpleMessageStream(ctx context.Context, text string) (<-chan string, error) {
	if text == "" {
		return nil, errors.New("text is empty")
	}

	r := make(<-chan ResponseMessageStream)

	var err error
	r, err = c.CreateMessageStream(ctx, &RequestMessageContent{
		Model:     DefaultModel,
		MaxTokens: DefaultMaxTokens,
		Messages: []*Message{
			{
				Role:    RoleUser,
				Content: text,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	ch := make(chan string)
	go func() {
		defer close(ch)
		var completed bool
		for completed == false {
			select {
			case rr, ok := <-r:
				if ok {
					if rr.Delta.Text != "" {
						ch <- rr.Delta.Text
					}
				} else {
					completed = true
				}
			case <-time.After(time.Second * 30): // 30s timeout
				completed = true
			}
		}
	}()
	return ch, nil
}
