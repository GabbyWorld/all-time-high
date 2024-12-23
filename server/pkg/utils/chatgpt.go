package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ChatGPTRequest struct {
	Model     string `json:"model"`
	Prompt    string `json:"prompt"`
	MaxTokens int    `json:"max_tokens"`
}

type ChatGPTResponse struct {
	Choices []struct {
		Text string `json:"text"`
	} `json:"choices"`
}

// GenerateDescription 使用真实的 OpenAI API
func GenerateDescription(apiKey, endpoint, name, prompt string) (string, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	requestBody := map[string]interface{}{
		"model": "gpt-4o", // 或 "gpt-4"（如果有权限）
		"messages": []map[string]string{
			{
				"role": "system",
				"content": `You're a creative storyteller and game designer with a talent for crafting engaging character descriptions in a vibrant gaming universe. Your task is to write a short, euphemistic description for an Agent in a player-vs-player battle arena.
										The description should subtly reflect the Agent's prompt without directly revealing its purpose or abilities, captivating players and sparking their imagination. Avoid mentioning the Agent's name in the description. Keep it under 160 characters.`,
			},
			{
				"role": "user",
				"content": fmt.Sprintf(`Agent's Name is: %s
																Agent's Prompt is: %s
																Description:`, name, prompt),
			},
		},
		"max_tokens": 1000, // todo: 需要根据实际情况调整
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// 解析响应
	var chatResp struct {
		ID      string `json:"id"`
		Object  string `json:"object"`
		Created int    `json:"created"`
		Choices []struct {
			Index   int `json:"index"`
			Message struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"message"`
			FinishReason string `json:"finish_reason"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(bodyBytes, &chatResp); err != nil {
		return "", err
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no description generated")
	}

	return chatResp.Choices[0].Message.Content, nil
}

// GenerateBattleOutcome 评估玩家对战的结果
func GenerateBattleOutcome(apiKey, endpoint, attName, attPrompt, defName, defPrompt string) (string, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	requestBody := map[string]interface{}{
		"model": "gpt-4o", // 或 "gpt-4"（如果有权限）
		"messages": []map[string]string{
			{
				"role": "system",
				"content": fmt.Sprintf(`You're a game system tasked with determining the outcome of battles in a player-vs-player arena featuring user-generated AI agents. Your role is to evaluate agents fairly and impartially, based only on the provided prompts, ensuring outcomes reflect their described abilities and how they might interact in an encounter.
																Analyze the following agents:
																- Attacker Agent Name: %s
																- Attacker Agent Prompt: %s
																- Defender Agent Name: %s
																- Defender Agent Prompt: %s
																1. Begin by stating the Attack Outcome:
																	- “Total Victory!” if the Attacker's abilities significantly outmatch the Defender's.
																	- “Narrow Victory!” if the Attacker has a slight edge.
																	- “Narrow Defeat!” if the Defender has a slight edge.
																	- “Crushing Defeat!” if the Defender's abilities significantly outmatch the Attacker's.
																2. Craft a story under 280 characters, reflecting the battle and its outcome.
																	- Mention both agents' names to make the narrative engaging.
																	- Avoid directly describing their abilities; focus on the imaginative depiction of how the battle unfolded.
																	- Ensure the story aligns with the logical implications of their abilities interacting (without bias from agent names).
																`, attName, attPrompt, defName, defPrompt),
			},
		},
		"max_tokens": 1000, // todo: 需要根据实际情况调整
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// 解析响应
	var chatResp struct {
		ID      string `json:"id"`
		Object  string `json:"object"`
		Created int    `json:"created"`
		Choices []struct {
			Index   int `json:"index"`
			Message struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"message"`
			FinishReason string `json:"finish_reason"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(bodyBytes, &chatResp); err != nil {
		return "", err
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no battle outcome generated")
	}

	return chatResp.Choices[0].Message.Content, nil
}
