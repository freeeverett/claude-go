package claude

import (
	"context"
	"fmt"
	"testing"
	"time"
)

var testKey = ""

func TestRequest(t *testing.T) {
	c := New(testKey)
	in := &RequestMessageContent{
		Model:     DefaultModel,
		MaxTokens: DefaultMaxTokens,
		Messages: []*Message{
			{
				Role:    RoleUser,
				Content: "Tell a joke",
			},
		},
	}
	r, err := c.CreateMessage(context.TODO(), in)
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
	r, err := c.CreateSimpleMessage(context.TODO(), text)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(r)
}
func TestRequestStream(t *testing.T) {

	c := New(testKey)
	in := &RequestMessageContent{
		Model:     DefaultModel,
		MaxTokens: DefaultMaxTokens,
		Messages: []*Message{
			{
				Role:    RoleUser,
				Content: "Tell a joke",
			},
		},
	}
	r := make(<-chan ResponseMessageStream)
	var err error
	r, err = c.CreateMessageStream(context.Background(), in)
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
	r, err = c.CreateSimpleMessageStream(context.Background(), text)
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
func TestMediaRequest(t *testing.T) {

	c := New(testKey)
	mediaData1 := &MessageContent{}
	mediaData1.Source.Data = "" // base64 encoded image data
	mediaData1.Source.MediaType = MediaTypeJPEG
	mediaData1.Source.Type = Base64ContentSourceType
	mediaData1.Type = ImageContentType

	mediaData2 := &MessageContent{}
	mediaData2.Text = "What is the artistic conception of this picture"
	mediaData2.Type = TextContentType

	mediaData := []*MessageContent{mediaData1, mediaData2}

	in := &RequestMessageContent{
		Model:     DefaultModel,
		MaxTokens: DefaultMaxTokens,
		Messages: []*Message{
			{
				Role:    RoleUser,
				Content: MarshalMediaContent(mediaData),
			},
		},
	}
	r, err := c.CreateMessage(context.TODO(), in)
	if err != nil {
		t.Error(err)
		return
	}
	if r != nil && len(r.Content) > 0 && r.Content[0].Type == "text" {
		t.Log(r.Content[0].Text)
		t.Log(r.Role, r.Model, r.Usage.InputTokens, r.Usage.OutputTokens)
	}
}
