package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/GabbyWorld/all-time-high-backend/internal/config"
	"github.com/GabbyWorld/all-time-high-backend/internal/logger"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	confirm "github.com/gagliardetto/solana-go/rpc/sendAndConfirmTransaction"
	"github.com/gagliardetto/solana-go/rpc/ws"
	"go.uber.org/zap"
)

// MetadataResponse 定义IPFS响应结构
type MetadataResponse struct {
	Metadata struct {
		Name        string `json:"name"`
		Symbol      string `json:"symbol"`
		Description string `json:"description"`
	} `json:"metadata"`
	MetadataUri string `json:"metadataUri"`
}

// CreateToken 创建Token并返回签名
func CreateToken(cfg *config.Config, imageURL string) (solana.Signature, error) {
	// 设置RPC端点
	clientRPC := rpc.New(cfg.Solana.RPCEndpoint)
	wsClient, err := ws.Connect(context.TODO(), cfg.Solana.WSRPCEndpoint)
	if err != nil {
		logger.Logger.Error("CreateToken: failed to connect to WebSocket RPC", zap.Error(err))
		return solana.Signature{}, fmt.Errorf("failed to connect to WebSocket RPC: %w", err)
	}
	defer wsClient.Close()

	// 解码签名者的私钥
	signerPrivateKey, err := solana.PrivateKeyFromBase58(cfg.Solana.SignerPrivateKey)
	if err != nil {
		logger.Logger.Error("CreateToken: invalid signer private key", zap.Error(err))
		return solana.Signature{}, fmt.Errorf("invalid signer private key: %w", err)
	}

	// 生成mintKeypair
	mintKeypair, err := solana.NewRandomPrivateKey()
	if err != nil {
		logger.Logger.Error("CreateToken: failed to generate mint keypair", zap.Error(err))
		return solana.Signature{}, fmt.Errorf("failed to generate mint keypair: %w", err)
	}
	mintPublicKey := mintKeypair.PublicKey()

	// 获取图片通过URL
	imageResp, err := http.Get(imageURL)
	if err != nil {
		logger.Logger.Error("CreateToken: failed to get image from URL", zap.Error(err))
		return solana.Signature{}, fmt.Errorf("failed to get image from URL: %w", err)
	}
	defer imageResp.Body.Close()

	if imageResp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(imageResp.Body)
		logger.Logger.Error("CreateToken: failed to fetch image", zap.String("status", imageResp.Status), zap.String("body", string(bodyBytes)))
		return solana.Signature{}, fmt.Errorf("failed to fetch image, status: %s", imageResp.Status)
	}

	imageData, err := io.ReadAll(imageResp.Body)
	if err != nil {
		logger.Logger.Error("CreateToken: failed to read image data", zap.Error(err))
		return solana.Signature{}, fmt.Errorf("failed to read image data: %w", err)
	}

	// 创建multipart form
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// 添加文件部分
	part, err := writer.CreateFormFile("file", "token.png")
	if err != nil {
		logger.Logger.Error("CreateToken: failed to create form file", zap.Error(err))
		return solana.Signature{}, fmt.Errorf("failed to create form file: %w", err)
	}
	_, err = part.Write(imageData)
	if err != nil {
		logger.Logger.Error("CreateToken: failed to write image data to form file", zap.Error(err))
		return solana.Signature{}, fmt.Errorf("failed to write image data to form file: %w", err)
	}

	// 添加其他表单字段
	fields := map[string]string{
		"name":        "overdosed2",
		"symbol":      "OVD2",
		"description": "This is an example token created via PumpPortal.fun",
		"twitter":     "",
		"telegram":    "",
		"website":     "",
		"showName":    "true",
	}

	for key, val := range fields {
		if err := writer.WriteField(key, val); err != nil {
			logger.Logger.Error("CreateToken: failed to write form field", zap.String("field", key), zap.Error(err))
			return solana.Signature{}, fmt.Errorf("failed to write form field %s: %w", key, err)
		}
	}

	if err := writer.Close(); err != nil {
		logger.Logger.Error("CreateToken: failed to close multipart writer", zap.Error(err))
		return solana.Signature{}, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	// 发送IPFS元数据存储请求
	req, err := http.NewRequest("POST", cfg.Solana.IPFSURL, &requestBody)
	if err != nil {
		logger.Logger.Error("CreateToken: failed to create IPFS request", zap.Error(err))
		return solana.Signature{}, fmt.Errorf("failed to create IPFS request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	clientHTTP := &http.Client{}
	ipfsResp, err := clientHTTP.Do(req)
	if err != nil {
		logger.Logger.Error("CreateToken: failed to send IPFS request", zap.Error(err))
		return solana.Signature{}, fmt.Errorf("failed to send IPFS request: %w", err)
	}
	defer ipfsResp.Body.Close()

	if ipfsResp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(ipfsResp.Body)
		logger.Logger.Error("CreateToken: IPFS request failed", zap.String("status", ipfsResp.Status), zap.String("body", string(bodyBytes)))
		return solana.Signature{}, fmt.Errorf("IPFS request failed, status: %s", ipfsResp.Status)
	}

	var metadataResp MetadataResponse
	if err := json.NewDecoder(ipfsResp.Body).Decode(&metadataResp); err != nil {
		logger.Logger.Error("CreateToken: failed to decode IPFS response", zap.Error(err))
		return solana.Signature{}, fmt.Errorf("failed to decode IPFS response: %w", err)
	}

	// 创建交易请求
	tradePayload := map[string]interface{}{
		"publicKey": signerPrivateKey.PublicKey().String(),
		"action":    "create",
		"tokenMetadata": map[string]string{
			"name":   metadataResp.Metadata.Name,
			"symbol": metadataResp.Metadata.Symbol,
			"uri":    metadataResp.MetadataUri,
		},
		"mint":             mintPublicKey.String(),
		"denominatedInSol": "true",
		"amount":           0,
		"slippage":         1,
		"priorityFee":      0,
		"pool":             "pump",
	}

	tradeBody, err := json.Marshal(tradePayload)
	if err != nil {
		logger.Logger.Error("CreateToken: failed to marshal trade payload", zap.Error(err))
		return solana.Signature{}, fmt.Errorf("failed to marshal trade payload: %w", err)
	}

	tradeReq, err := http.NewRequest("POST", cfg.Solana.TradeURL, bytes.NewBuffer(tradeBody))
	if err != nil {
		logger.Logger.Error("CreateToken: failed to create trade request", zap.Error(err))
		return solana.Signature{}, fmt.Errorf("failed to create trade request: %w", err)
	}
	tradeReq.Header.Set("Content-Type", "application/json")

	tradeResp, err := clientHTTP.Do(tradeReq)
	if err != nil {
		logger.Logger.Error("CreateToken: failed to send trade request", zap.Error(err))
		return solana.Signature{}, fmt.Errorf("failed to send trade request: %w", err)
	}
	defer tradeResp.Body.Close()

	// 读取响应体
	bodyBytes, err := io.ReadAll(tradeResp.Body)
	if err != nil {
		logger.Logger.Error("CreateToken: failed to read trade response body", zap.Error(err))
		return solana.Signature{}, fmt.Errorf("failed to read trade response body: %w", err)
	}

	// 输出响应内容类型和长度
	contentType := tradeResp.Header.Get("Content-Type")
	logger.Logger.Info("CreateToken: trade response received", zap.String("content_type", contentType), zap.Int("body_length", len(bodyBytes)))

	// 检查响应是否成功
	if tradeResp.StatusCode != http.StatusOK {
		logger.Logger.Error("CreateToken: trade request failed", zap.String("status", tradeResp.Status), zap.String("body", string(bodyBytes)))
		return solana.Signature{}, fmt.Errorf("trade request failed, status: %s", tradeResp.Status)
	}

	// 将响应体解析为交易对象
	tx, err := solana.TransactionFromBytes(bodyBytes)
	if err != nil {
		logger.Logger.Error("CreateToken: failed to parse transaction from response", zap.Error(err))
		return solana.Signature{}, fmt.Errorf("failed to parse transaction from response: %w", err)
	}

	// 获取最新的blockhash
	recentBlockhashResp, err := clientRPC.GetLatestBlockhash(context.TODO(), rpc.CommitmentFinalized)
	if err != nil {
		logger.Logger.Error("CreateToken: failed to get latest blockhash", zap.Error(err))
		return solana.Signature{}, fmt.Errorf("failed to get latest blockhash: %w", err)
	}
	tx.Message.RecentBlockhash = recentBlockhashResp.Value.Blockhash

	// 签名交易
	_, err = tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if signerPrivateKey.PublicKey().Equals(key) {
			return &signerPrivateKey
		} else if mintKeypair.PublicKey().Equals(key) {
			return &mintKeypair
		}
		return nil
	})
	if err != nil {
		logger.Logger.Error("CreateToken: failed to sign transaction", zap.Error(err))
		return solana.Signature{}, fmt.Errorf("failed to sign transaction: %w", err)
	}

	// 发送并确认交易
	sig, err := confirm.SendAndConfirmTransactionWithTimeout(
		context.TODO(),
		clientRPC,
		wsClient,
		tx,
		time.Second*90,
	)
	if err != nil {
		logger.Logger.Error("CreateToken: failed to send and confirm transaction", zap.Error(err))
		return solana.Signature{}, fmt.Errorf("failed to send and confirm transaction: %w", err)
	}

	return sig, nil
}

// GetTokenAddress 根据签名获取Token地址
func GetTokenAddress(cfg *config.Config, sig solana.Signature) (string, error) {
	// 设置 Solana RPC 客户端（可以选择主网或测试网）
	client := rpc.New(cfg.Solana.RPCEndpoint)

	// 获取交易详情
	maxVersion := uint64(0)
	txInfo, err := client.GetParsedTransaction(
		context.TODO(),
		sig,
		&rpc.GetParsedTransactionOpts{
			MaxSupportedTransactionVersion: &maxVersion,
		},
	)
	if err != nil {
		logger.Logger.Error("GetTokenAddress: failed to get parsed transaction", zap.Error(err))
		return "", fmt.Errorf("failed to get parsed transaction: %w", err)
	}

	if txInfo == nil {
		logger.Logger.Error("GetTokenAddress: transaction info is nil")
		return "", fmt.Errorf("transaction info is nil")
	}

	// 解析交易中的每个指令，查找 Token 地址
	tokenProgramID := solana.MustPublicKeyFromBase58(cfg.Solana.TokenProgramID)
	for _, ix := range txInfo.Transaction.Message.Instructions {
		if ix.ProgramId.Equals(tokenProgramID) {
			// Token 地址通常会出现在accounts数组的某一位置
			// 根据 Token 指令类型不同（创建账户、初始化账户等），我们可以提取 Token 地址
			if len(ix.Accounts) > 0 {
				logger.Logger.Info("GetTokenAddress: found token account address", zap.String("token_address", ix.Accounts[0].String()))
				return ix.Accounts[0].String(), nil
			}
		}
	}

	logger.Logger.Error("GetTokenAddress: no token account address found in transaction")
	return "", fmt.Errorf("no token account address found in transaction")
}
