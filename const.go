package claude

type Client struct {
	apiKey string
	domain string
}

type role string

const (
	RoleUser      role = "user"
	RoleAssistant role = "assistant"
)

//

type model string

const (
	ModelClaude3Opus20240229   model = "claude-3-opus-20240229"
	ModelClaude3Sonnet20240229 model = "claude-3-sonnet-20240229"
)
