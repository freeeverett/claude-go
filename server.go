package claude_go

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

func New(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		domain: "https://api.anthropic.com/",
	}
}

func (c *Client) SetDomain(domain string) {
	// SetDomain 设置请求的域名，默认为 https://api.anthropic.com/
	c.domain = domain
}

func (c *Client) CreateMessage(in *RequestMessageContent) (*ResponseMessageContent, error) {
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
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("content-type", "application/json")

	r, errDo := cli.Do(req)
	if errDo != nil {
		return nil, errDo
	}
	if r.StatusCode != 200 {
		var errRes = ResponseError{}
		_ = json.NewDecoder(r.Body).Decode(&errRes)
		errRes.Code = r.StatusCode
		return nil, &errRes
	}
	var res ResponseMessageContent
	return &res, json.NewDecoder(r.Body).Decode(&res)
}
func (c *Client) CreateSimpleMessage(text string) (string, error) {
	if text == "" {
		return "", errors.New("text is empty")
	}
	r, err := c.CreateMessage(&RequestMessageContent{
		Model:     ModelClaude3Sonnet20240229,
		MaxTokens: 4096,
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
func (c *Client) CreateMessageStream(in *RequestMessageContent) (<-chan ResponseMessageStream, error) {
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
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("anthropic-beta", "messages-2023-12-15")

	r, errDo := cli.Do(req)
	if errDo != nil {
		return nil, errDo
	}
	if r.StatusCode != 200 {
		var errRes = ResponseError{}
		_ = json.NewDecoder(r.Body).Decode(&errRes)
		errRes.Code = r.StatusCode
		return nil, &errRes
	}
	reader := bufio.NewReader(r.Body)
	//
	ch := make(chan ResponseMessageStream)
	go func() {
		defer close(ch)
		for {
			line, err1 := reader.ReadBytes('\n')
			if err1 != nil {
				break
			}
			if len(line) > 6 && string(line[:5]) == "data:" {
				data := string(line[6:])
				var res ResponseMessageStream
				_ = json.Unmarshal([]byte(data), &res)
				ch <- res
			}
			//
		}
	}()
	return ch, nil
}
func (c *Client) CreateSimpleMessageStream(text string) (<-chan string, error) {
	if text == "" {
		return nil, errors.New("text is empty")
	}

	r := make(<-chan ResponseMessageStream)

	var err error
	r, err = c.CreateMessageStream(&RequestMessageContent{
		Model:     ModelClaude3Sonnet20240229,
		MaxTokens: 4096,
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
		for {
			exit := false
			select {
			case rr, ok := <-r:
				if ok {
					if rr.Delta.Text != "" {
						ch <- rr.Delta.Text
					}
				} else {
					exit = true
				}
			case <-time.After(time.Second * 30): // 30s timeout
				exit = true
			}
			if exit {
				break
			}
		}
	}()
	return ch, nil
}
