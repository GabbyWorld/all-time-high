package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"strings"
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

	// 将结果四舍五入到小数点后 7 位
	priceVal = math.Round(priceVal*1e7) / 1e7

	return priceVal, nil
}

// GetTokenPrice 根据给定的 tokenAddress 从外部服务获取市值（价格）
func GetTokenVsSOLPrice(tokenAddress string) (float64, error) {
	if tokenAddress == "" {
		return 0, errors.New("tokenAddress is empty")
	}

	resp, err := http.Get(fmt.Sprintf("https://api.jup.ag/price/v2?ids=%s&vsToken=So11111111111111111111111111111111111111112", tokenAddress))
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

	// 将结果四舍五入到小数点后 10 位
	priceVal = math.Round(priceVal*1e10) / 1e10

	return priceVal, nil
}

// GetMultipleTokenPrice 根据给定的 tokenAddresses 从外部服务获取多个代币的市值（价格）
func GetMultipleTokenPrice(tokenAddresses []string) (map[string]float64, error) {
	if len(tokenAddresses) == 0 {
		return nil, errors.New("tokenAddresses is empty")
	}

	resp, err := http.Get(fmt.Sprintf("https://api.jup.ag/price/v2?ids=%s", strings.Join(tokenAddresses, ",")))
	if err != nil {
		log.Printf("Failed to send GET request: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		return nil, err
	}

	var apiResponse ApiResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		log.Printf("Failed to parse JSON: %v", err)
		return nil, err
	}

	prices := make(map[string]float64)
	for _, tokenAddress := range tokenAddresses {
		priceStr := apiResponse.Data[tokenAddress].Price
		priceVal, err := priceStr.Float64()
		if err != nil {
			log.Printf("Failed to convert priceStr to float64 for token %s: %v", tokenAddress, err)
			return nil, err
		}

		// 将结果四舍五入到小数点后 7 位
		priceVal = math.Round(priceVal*1e7) / 1e7

		prices[tokenAddress] = priceVal
	}

	return prices, nil
}

// GetMultipleTokenVsSOLPrice 根据给定的 tokenAddresses 从外部服务获取多个代币相对于 SOL 的价格
func GetMultipleTokenVsSOLPrice(tokenAddresses []string) (map[string]float64, error) {
	if len(tokenAddresses) == 0 {
		return nil, errors.New("tokenAddresses is empty")
	}

	url := fmt.Sprintf(
		"https://api.jup.ag/price/v2?ids=%s&vsToken=So11111111111111111111111111111111111111112",
		strings.Join(tokenAddresses, ","),
	)

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to send GET request: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		return nil, err
	}

	var apiResponse ApiResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		log.Printf("Failed to parse JSON: %v", err)
		return nil, err
	}

	prices := make(map[string]float64)
	for _, tokenAddress := range tokenAddresses {
		priceStr := apiResponse.Data[tokenAddress].Price
		priceVal, err := priceStr.Float64()
		if err != nil {
			log.Printf("Failed to convert priceStr to float64 for token %s: %v", tokenAddress, err)
			return nil, err
		}

		// 将结果四舍五入到小数点后 10 位
		priceVal = math.Round(priceVal*1e10) / 1e10

		prices[tokenAddress] = priceVal
	}

	return prices, nil
}
