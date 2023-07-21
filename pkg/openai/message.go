package openai

// Message is a single message in the OpenAI conversation
type Message struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}

func NewUserMessage(content string) Message {
	return Message{
		Content: content,
		Role:    "user",
	}
}

func NewSystemMessage(content string) Message {
	return Message{
		Content: content,
		Role:    "system",
	}
}

// Conversation is a conversation between the user and the bot
type Conversation struct {
	Messages []Message
}

// NewConversation creates a new conversation
func NewConversation(systemMessage Message, examples ...Message) *Conversation {
	c := Conversation{}
	c.Messages = append(c.Messages, systemMessage)

	for _, example := range examples {
		c.Messages = append(c.Messages, example)
	}

	return &c
}

// AddMessage adds a message to the conversation from a response
func (c *Conversation) AddMessage(response OpenAiCompletionResponse) {
	for _, choice := range response.Choices {
		c.Messages = append(c.Messages, choice.Message)
	}
}

// OpenAiCompletionRequest is the request body to the OpenAI API
type OpenAiCompletionRequest struct {
	Messages []Message `json:"messages"`
}

// OpenAiCompletionResponse is the response body from the OpenAI API
type OpenAiCompletionResponse struct {
	Choices []struct {
		FinishReason string  `json:"finish_reason"`
		Index        int     `json:"index"`
		Message      Message `json:"message"`
	} `json:"choices"`
	Created int    `json:"created"`
	Id      string `json:"id"`
	Model   string `json:"model"`
	Object  string `json:"object"`
	Usage   struct {
		CompletionTokens int `json:"completion_tokens"`
		PromptTokens     int `json:"prompt_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}