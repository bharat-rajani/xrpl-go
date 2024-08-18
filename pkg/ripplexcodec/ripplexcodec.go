package ripplexcodec

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/xrpscan/xrpl-go/pkg/base58"
	"github.com/xrpscan/xrpl-go/pkg/utils"
)

const (
	// Lengths in bytes
	ClassicAddressLength   = 20
	AccountPublicKeyLength = 33
	FamilySeedLength       = 16
	NodePublicKeyLength    = 33

	// Account/classic address prefix - value is 0
	ClassicAddressPrefix = 0x00
	// Account public key prefix - value is 35
	AccountPublicKeyPrefix = 0x23
	// Family seed prefix - value is 33
	FamilySeedPrefix = 0x21
	// Node/validation public key prefix - value is 28
	NodePublicKeyPrefix = 0x1C
	// ED25519 prefix - value is 237
	ED25519Prefix = 0xED
)

type PrefixBytes struct {
	MainNet, TestNet []byte
}

func GetPrefixBytes() PrefixBytes {
	return PrefixBytes{
		MainNet: []byte{0x05, 0x44},
		TestNet: []byte{0x04, 0x93},
	}
}

// XAddress
type XAddress struct {
	AccountID []byte
	Tag       *int32
	IsTest    bool
}

// IsValidXAddress
func IsValidXAddress(xAddress string) bool {
	_, err := DecodeXAddress(xAddress)
	return err == nil
}

func DecodeXAddress(xAddress string) (*XAddress, error) {
	decoded, err := decodeChecked(xAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to decode xAddress: %w", err)
	}
	accountId := decoded[2:22]
	isTest, err := isTestXAddress(decoded)
	if err != nil {
		return nil, fmt.Errorf("failed to decode xAddress: %w", err)
	}
	tag, err := TagFromXAddress(decoded)
	if err != nil {
		return nil, fmt.Errorf("failed to decode xAddress tag: %w", err)
	}
	return &XAddress{
		AccountID: accountId,
		Tag:       tag,
		IsTest:    isTest,
	}, nil
}

func decodeChecked(base58str string) ([]byte, error) {
	decoded := base58.Decode(base58str)
	if len(decoded) < 5 {
		return nil, errors.New("invalid input size: decoded data must have length >= 5")
	}
	if !verifyCheckSum(decoded) {
		return nil, errors.New("checksum invalid")
	}
	return decoded[:len(decoded)-4], nil
}

func verifyCheckSum(xAddrDecoded []byte) bool {
	h := sha256.Sum256(xAddrDecoded[:len(xAddrDecoded)-4])
	h2 := sha256.Sum256(h[:])
	computed := h2[:4]
	checkSum := xAddrDecoded[len(xAddrDecoded)-4:]
	return bytes.Equal(checkSum, computed)
}

func isTestXAddress(xAddressDecoded []byte) (bool, error) {
	decodedPrefix := xAddressDecoded[:2]

	if bytes.Equal(decodedPrefix, GetPrefixBytes().MainNet) {
		return false, nil
	}
	if bytes.Equal(decodedPrefix, GetPrefixBytes().TestNet) {
		return true, nil
	}
	return false, errors.New("invalid X-address: bad prefix")
}

// TagFromXAddress returns the destination tag extracted from the suffix of the X-Address.
//
//	Args:
//	    buffer: The buffer to extract a destination tag from.
//
//	Returns:
//	    The destination tag extracted from the suffix of the X-Address.
func TagFromXAddress(xAddrDecoded []byte) (*int32, error) {

	flag := xAddrDecoded[22]
	if flag >= 2 {
		return nil, errors.New("unsupported X-Address")
	}
	if flag == 1 {
		t := binary.LittleEndian.Uint32(xAddrDecoded[23:27])
		return utils.PointerOf(int32(t)), nil
	}
	if flag != 0 {
		return nil, errors.New("flag must be zero to indicate no tag")
	}
	zeroBytes, _ := hex.DecodeString("0000000000000000")
	if !bytes.Equal(zeroBytes, xAddrDecoded[23:23+8]) {
		return nil, errors.New("remaining bytes must be zero")
	}
	return nil, nil
}

func XAddressToClassicAddress(xAddress string) (string, *int32, bool) {
	xAddr, err := DecodeXAddress(xAddress)
	if err != nil {
		return "", nil, false
	}
	classicAddrBytes := xAddr.AccountID
	address, err := EncodeClassicAddress(classicAddrBytes)
	if err != nil {
		return "", nil, xAddr.IsTest
	}
	return address, xAddr.Tag, xAddr.IsTest
}

func EncodeClassicAddress(classicAddr []byte) (string, error) {

	encodedClassicAddr, err := EncodeAddress(classicAddr, ClassicAddressPrefix, ClassicAddressLength)
	if err != nil {
		return "", err
	}
	return encodedClassicAddr, nil
}

func EncodeAddress(input []byte, prefix byte, expectedLen int) (string, error) {
	if len(input) != expectedLen {
		return "", errors.New("unexpected payload length: len(input) does not match expected_length")
	}

	return base58.CheckEncode(input, prefix), nil
}
