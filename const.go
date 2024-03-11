package claude

type Client struct {
	apiKey string
	domain string
}

type Role string

const (
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

//

const (
	ModelClaude3Opus20240229   Model = "claude-3-opus-20240229"
	ModelClaude3Sonnet20240229 Model = "claude-3-sonnet-20240229"
)
