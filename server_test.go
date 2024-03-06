package claude

import (
	"fmt"
	"testing"
	"time"
)

var testKey = ""

func TestRequest(t *testing.T) {
	c := New(testKey)
	in := &RequestMessageContent{
		Model:     ModelClaude3Sonnet20240229,
		MaxTokens: 4096,
		Messages: []*Message{
			{
				Role:    RoleUser,
				Content: "Tell a joke",
			},
		},
	}
	r, err := c.CreateMessage(in)
	if err != nil {
		t.Error(err)
		return
	}
	if r != nil && len(r.Content) > 0 && r.Content[0].Type == "text" {
		t.Log(r.Content[0].Text)
		t.Log(r.Role, r.Model, r.Usage.InputTokens, r.Usage.OutputTokens)
	}
}
func TestSimpleRequest(t *testing.T) {
	c := New(testKey)
	text := "Tell a joke"
	r, err := c.CreateSimpleMessage(text)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(r)
}
func TestRequestStream(t *testing.T) {

	c := New(testKey)
	in := &RequestMessageContent{
		Model:     ModelClaude3Sonnet20240229,
		MaxTokens: 4096,
		Messages: []*Message{
			{
				Role:    RoleUser,
				Content: "Tell a joke",
			},
		},
	}
	r := make(<-chan ResponseMessageStream)
	var err error
	r, err = c.CreateMessageStream(in)
	if err != nil {
		t.Error(err)
		return
	}
	for {
		exit := false
		select {
		case rr, ok := <-r:
			if ok {
				if rr.Delta.Text != "" {
					fmt.Println(rr.Delta.Text)
				}
				if len(rr.Message.Content) > 0 {
					fmt.Println(rr.Message.Content[0].Text)
				}
			} else {
				fmt.Println("channel closed")
				exit = true
			}
		case <-time.After(time.Second * 100):
			fmt.Println("timeout")
			exit = true
		}
		if exit {
			break
		}
	}
}
func TestSimpleRequestStream(t *testing.T) {

	c := New(testKey)
	text := "Tell a joke"
	r := make(<-chan string)
	var err error
	r, err = c.CreateSimpleMessageStream(text)
	if err != nil {
		t.Error(err)
		return
	}
	for {
		exit := false
		select {
		case rr, ok := <-r:
			if ok {
				fmt.Println(rr)
			} else {
				fmt.Println("channel closed")
				exit = true
			}
		case <-time.After(time.Second * 300):
			fmt.Println("timeout")
			exit = true
		}
		if exit {
			break
		}
	}
}
