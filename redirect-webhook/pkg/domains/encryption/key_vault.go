package encryption

import (
	"bytes"
	"conductor.com.br/pier/hooks/pkg/conf/global"
	"conductor.com.br/pier/hooks/pkg/helpers"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"strings"
	"time"
)

type KeyVaultStatus string

func (status KeyVaultStatus) Value() string {
	return string(status)
}

type KeyVault struct {
	helpers.GormHooks
	ID                      string
	IssuerPublicKey         []byte
	AESKey                  []byte
	Enabled                 bool
	IssuerKeyExpirationDate time.Time
	AESKeyExpirationDate    time.Time
	IssuerId                int

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (k *KeyVault) SetInitialStatusAndExpiration() {
	k.Enabled = true
	k.AESKeyExpirationDate = time.Now().AddDate(0, global.Viper.GetInt("encryption.aesKey.monthsToExpire"), 0)
}

func (k KeyVault) IssuerKeyHasExpired() bool {
	return time.Now().After(k.IssuerKeyExpirationDate)
}

func (k KeyVault) GetRSAPublicKey() (*rsa.PublicKey, error) {
	block, _ := pem.Decode(k.IssuerPublicKey)
	if block == nil {
		return nil, errors.New("public key is not valid. Must start with -----BEGIN  and end with -----END *-----")
	}
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			return nil, err
		}
	}

	ifc, err := x509.ParsePKIXPublicKey(b)
	if err != nil {
		return nil, err
	}

	key, ok := ifc.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("could not parse to RSA PublicKey")
	}

	return key, nil
}

func (k KeyVault) EncryptWithPublicKey(value []byte) (string, error) {
	publicKey, err := k.GetRSAPublicKey()
	if err != nil {
		return "", err
	}

	encryptedMsg, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, value)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encryptedMsg), nil
}

func (k KeyVault) EncryptWithPublicKeyMock(value []byte) (string, error) {
	publicKey, err := k.GetRSAPublicKey()
	if err != nil {
		return "", err
	}
	pubKey, _ := k.GetRSAPublicKey()

	byteData := make([]byte, pubKey.Size())
	for i := range byteData {
		byteData[i] = 1
	}

	encryptedMsg, err := rsa.EncryptPKCS1v15(bytes.NewReader(byteData), publicKey, value)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encryptedMsg), nil
}

func (k KeyVault) EncryptAesWithPublicKey() (string, error) {
	return k.EncryptWithPublicKey(k.AESKey)
}

func (k KeyVault) EncryptAesWithPublicKeyMock() (string, error) {
	return k.EncryptWithPublicKeyMock(k.AESKey)
}

// TODO improve this function
func (k KeyVault) ValidatePublicKey() error {
	block, _ := pem.Decode(k.IssuerPublicKey)
	if block == nil {
		return errors.New("public key is not valid. Must start with -----BEGIN  and end with -----END *-----, and each line break should contain a '\\n'")
	}
	return nil
}

func (k *KeyVault) CreateAesKey() error {
	// Generate a 256bit/32byte array
	genkey := make([]byte, 32)

	// Randomized each byte's value
	_, err := rand.Read(genkey)

	k.AESKey = genkey

	return err
}

func (k *KeyVault) CreateAesKeyMock() error {
	defaultAESKey := []byte{249, 204, 3, 19, 252, 165, 224, 60, 71, 89, 253, 68, 12, 234, 92, 180, 28, 251, 235, 180, 80, 6, 241, 9, 79, 112, 216, 229, 106, 210, 186, 153}

	k.AESKey = defaultAESKey

	return nil
}

func (k KeyVault) EncryptWithAes(text string) (string, error) {
	block, err := aes.NewCipher(k.AESKey)
	if err != nil {
		return "", err
	}

	msg := Pad([]byte(text))
	byteArrayOutput := make([]byte, len(msg))

	cbc := cipher.NewCBCEncrypter(block, make([]byte, aes.BlockSize))
	cbc.CryptBlocks(byteArrayOutput, msg)
	finalMsg := removeBase64Padding(base64.StdEncoding.EncodeToString(byteArrayOutput))
	return finalMsg, nil
}

func (k KeyVault) DecryptWithAes(text string) (string, error) {
	block, err := aes.NewCipher(k.AESKey)
	if err != nil {
		return "", err
	}

	decodedMsg, e := base64.StdEncoding.DecodeString(addBase64Padding(text))
	if e != nil {
		return "", err
	}

	if (len(decodedMsg) % aes.BlockSize) != 0 {
		return "", errors.New("blocksize must be multipe of decoded message length")
	}

	msg := decodedMsg

	decrypter := cipher.NewCBCDecrypter(block, make([]byte, aes.BlockSize))

	decrypter.CryptBlocks(msg, msg)

	unpadMsg, err := Unpad(msg)
	if err != nil {
		return "", err
	}

	return string(unpadMsg), nil
}

// Adds "=" character to a string
func addBase64Padding(value string) string {
	m := len(value) % 4
	if m != 0 {
		value += strings.Repeat("=", 4-m)
	}

	return value
}

// Removes "=" character from a string
func removeBase64Padding(value string) string {
	return strings.Replace(value, "=", "", -1)
}

// Adds "=" character to a byte array
func Pad(src []byte) []byte {
	padding := aes.BlockSize - len(src)%aes.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

// Removes "=" character from a byte array
func Unpad(src []byte) ([]byte, error) {
	length := len(src)
	unpadding := int(src[length-1])

	if unpadding > length {
		return nil, errors.New("unpad error. This could happen when incorrect encryption key is used")
	}

	return src[:(length - unpadding)], nil
}
