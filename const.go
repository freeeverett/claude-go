package claude

const (
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

//

const (
	Opus20240229   Model = "claude-3-opus-20240229"
	Sonnet20240229 Model = "claude-3-sonnet-20240229"
	Haiku20240307  Model = "claude-3-haiku-20240307"
)

const (
	DefaultModel     = Sonnet20240229
	DefaultMaxTokens = 4096
	DefaultDomain    = "https://api.anthropic.com/"
)

const (
	TextContentType         = "text"
	ImageContentType        = "image"
	Base64ContentSourceType = "base64"
	MediaTypeJPEG           = "image/jpeg"
)

const (
	anthropicVersion = "2023-06-01"
)
