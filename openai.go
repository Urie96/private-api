package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/alevinval/sse/pkg/decoder"
	"github.com/pkg/errors"
)

type OpenAIParams struct {
	Model            string   `json:"model"`
	Prompt           string   `json:"prompt" binding:"required"`
	Temperature      float32  `json:"temperature"`
	MaxTokens        int64    `json:"max_tokens"`
	TopP             float32  `json:"top_p"`
	FrequencyPenalty float32  `json:"frequency_penalty"`
	PresencePenalty  float32  `json:"presence_penalty"`
	Stop             []string `json:"stop"`
	Stream           bool     `json:"stream"`
}

func OpenAIRequest(ctx context.Context, p OpenAIParams) (io.Reader, error) {
	b, _ := json.Marshal(p)
	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/completions", bytes.NewReader(b))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	envs := GetEnvs()
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", envs.OpenAI.SecretKey))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	dec := decoder.New(resp.Body)
	type StreamEvent struct {
		Choices []struct {
			Text string `json:"text"`
		} `json:""choices`
	}
	reader, writer := io.Pipe()
	fmt.Println()
	go func() {
		for {
			event, err := dec.Decode()
			if err != nil {
				if err != io.EOF {
					fmt.Println("read error:", err)
				}
				writer.Close()
				break
			}
			e := &StreamEvent{}
			json.Unmarshal([]byte(event.GetData()), e)

			for _, c := range e.Choices {
				fmt.Print(c.Text)
				writer.Write([]byte(c.Text))
			}
		}
	}()

	return reader, nil
}
