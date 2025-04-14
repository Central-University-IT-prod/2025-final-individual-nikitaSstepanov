package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	e "github.com/nikitaSstepanov/tools/error"
)

type Ai struct {
	cfg *Config
}

type Config struct {
	OpenRouterAddress string `yaml:"address"`
	OpenRouterKey     string `env:"OPEN_ROUTER_KEY"`
	LLM               string `yaml:"llm"`
}

func New(cfg *Config) *Ai {
	return &Ai{
		cfg: cfg,
	}
}

func (a *Ai) GenText(prompt string) (string, e.Error) {
	body := &Body{
		Model: a.cfg.LLM,
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	var buffer bytes.Buffer

	err := json.NewEncoder(&buffer).Encode(body)
	if err != nil {
		return "", e.InternalErr.WithErr(err)
	}

	req, err := http.NewRequest("POST", a.cfg.OpenRouterAddress, &buffer)
	if err != nil {
		return "", e.InternalErr.WithErr(err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.cfg.OpenRouterKey))

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", e.InternalErr.WithErr(err)
	}

	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return "", e.InternalErr.WithErr(err)
	}

	var resp Response

	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		return "", e.InternalErr.WithErr(err)
	}

	if len(resp.Choices) == 0 {
		return "", e.BadInputErr
	}

	return resp.Choices[0].Message.Content, nil
}
