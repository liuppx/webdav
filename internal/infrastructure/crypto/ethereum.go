package crypto

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	ErrInvalidSignature       = errors.New("invalid signature")
	ErrInvalidSignatureLength = errors.New("invalid signature length")
	ErrSignatureMismatch      = errors.New("signature address mismatch")
)

// EthereumSigner 以太坊签名验证器
type EthereumSigner struct{}

// NewEthereumSigner 创建以太坊签名验证器
func NewEthereumSigner() *EthereumSigner {
	return &EthereumSigner{}
}

// VerifySignature 验证以太坊签名
func (s *EthereumSigner) VerifySignature(message, signatureHex, expectedAddress string) error {
	// 移除 0x 前缀
	signatureHex = strings.TrimPrefix(signatureHex, "0x")
	
	// 解码签名
	signature, err := hex.DecodeString(signatureHex)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidSignature, err)
	}
	
	if len(signature) != 65 {
		return ErrInvalidSignatureLength
	}
	
	// 调整 v 值（MetaMask 等钱包会加 27）
	if signature[64] >= 27 {
		signature[64] -= 27
	}
	
	// 构建以太坊签名消息
	hash := s.hashMessage(message)
	
	// 恢复公钥
	pubKey, err := crypto.SigToPub(hash.Bytes(), signature)
	if err != nil {
		return fmt.Errorf("failed to recover public key: %w", err)
	}
	
	// 从公钥生成地址
	recoveredAddress := crypto.PubkeyToAddress(*pubKey)
	
	// 比较地址
	if !strings.EqualFold(recoveredAddress.Hex(), expectedAddress) {
		return fmt.Errorf("%w: expected %s, got %s", 
			ErrSignatureMismatch, expectedAddress, recoveredAddress.Hex())
	}
	
	return nil
}

// hashMessage 哈希消息（以太坊签名消息格式）
func (s *EthereumSigner) hashMessage(message string) common.Hash {
	prefix := fmt.Sprintf("\x19Ethereum Signed Message:\n%d", len(message))
	return crypto.Keccak256Hash([]byte(prefix + message))
}

// IsValidAddress 验证地址格式
func (s *EthereumSigner) IsValidAddress(address string) bool {
	return common.IsHexAddress(address)
}

