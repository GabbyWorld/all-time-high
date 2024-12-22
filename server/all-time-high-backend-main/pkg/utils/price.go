package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

type ApiResponse struct {
	Data      map[string]TokenData `json:"data"`
	TimeTaken float64              `json:"timeTaken"`
}

type TokenData struct {
	ID    string      `json:"id"`
	Type  string      `json:"type"`
	Price json.Number `json:"price"`
}

// GetTokenPrice 根据给定的 tokenAddress 从外部服务获取市值（价格）
func GetTokenPrice(tokenAddress string) (float64, error) {
	if tokenAddress == "" {
		return 0, errors.New("tokenAddress is empty")
	}

	resp, err := http.Get(fmt.Sprintf("https://api.jup.ag/price/v2?ids=%s", tokenAddress))
	if err != nil {
		log.Printf("Failed to send GET request: %v", err)
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		return 0, err
	}

	var apiResponse ApiResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		log.Printf("Failed to parse JSON: %v", err)
		return 0, err
	}

	priceStr := apiResponse.Data[tokenAddress].Price
	priceVal, err := priceStr.Float64()
	if err != nil {
		log.Printf("Failed to convert gabbyPriceStr to float64: %v", err)
		return 0, err
	}

	return priceVal, nil
}
