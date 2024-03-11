package claude

// https://docs.anthropic.com/claude/reference/getting-started-with-the-api

type Message struct {
	Role    Role   `json:"role"`
	Content string `json:"content"`
}

type MessageContent struct {
	Type   string `json:"type"` // image //
	Source struct {
		Type      string `json:"type"`       // base64
		MediaType string `json:"media_type"` // image/png, image/jpeg, image/gif, image/webp
		Data      string `json:"data"`
	} `json:"source"`
	Text string `json:"text"`
}

type Model string

type RequestMessageContent struct {
	Model       Model      `json:"model"`
	MaxTokens   int64      `json:"max_tokens"`
	Messages    []*Message `json:"messages"`
	System      string     `json:"system,omitempty"` // system prompt
	Stream      bool       `json:"stream"`
	Temperature float64    `json:"temperature,omitempty"` // default 1.0  0.0 - 1.0
	TopP        float64    `json:"top_p,omitempty"`       // default 1.0  0.0 - 1.0 // default only use temperature
	TopK        int64      `json:"top_k,omitempty"`       // default only use temperature
}

type ResponseMessageContent struct {
	Id      string `json:"id"`
	Type    string `json:"type"`
	Role    string `json:"role"`
	Content []struct {
		Text string `json:"text"`
		Type string `json:"type"` // only support text
	} `json:"content"`
	Model        Model   `json:"model"`
	StopReason   *string `json:"stop_reason"` // [end_turn|max_tokens|stop_sequence]
	StopSequence *string `json:"stop_sequence"`
	Usage        struct {
		InputTokens  int64 `json:"input_tokens"`
		OutputTokens int64 `json:"output_tokens"`
	} `json:"usage"`
}

type ResponseMessageStream struct {
	Type         string                 `json:"type"` // message_start - content_block_start - ping - content_block_delta - content_block_stop - message_delta - message_stop
	Index        int64                  `json:"index"`
	Message      ResponseMessageContent `json:"message"`
	ContentBlock struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content_block"`
	Delta struct {
		Type         string  `json:"type"`
		Text         string  `json:"text"`
		StopReason   string  `json:"stop_reason"`
		StopSequence *string `json:"stop_sequence"`
	} `json:"delta"`
	Usage struct {
		OutputTokens int64 `json:"output_tokens"`
		InputTokens  int64 `json:"input_tokens"`
	} `json:"usage"`
}
