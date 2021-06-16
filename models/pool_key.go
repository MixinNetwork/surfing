package models

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/MixinNetwork/bot-api-go-client"
	"github.com/MixinNetwork/surfing/config"
	"github.com/MixinNetwork/surfing/session"
	"github.com/ugorji/go/codec"
)

const encryptionHeaderLength = 16

var poolKeysColumnsFull = []string{"user_id", "session_id", "session_key", "pin_token", "encrypted_pin", "encryption_header", "surfing_key", "created_at"}

func (k *PoolKey) valuesFull() []interface{} {
	return []interface{}{k.UserId, k.SessionId, k.SessionKey, k.PinToken, k.EncryptedPIN, k.EncryptionHeader, k.SurfingKey, k.CreatedAt}
}

type PoolKey struct {
	UserId           string
	SessionId        string
	SessionKey       string
	PinToken         string
	EncryptedPIN     string
	EncryptionHeader []byte
	SurfingKey       string
	CreatedAt        time.Time

	PlainPIN string
}

func GeneratePoolKey(ctx context.Context, mainnetAddress string) (*PoolKey, error) {
	pub, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	privateKeyString := hex.EncodeToString(privateKey[:])
	sessionSecret := base64.RawURLEncoding.EncodeToString(pub[:])
	user, err := bot.CreateUser(ctx, sessionSecret, fmt.Sprintf("Surfing %x", md5.Sum(pub[:])), config.ClientId, config.SessionId, config.SessionKey)
	if err != nil {
		return nil, err
	}
	key := &PoolKey{
		UserId:     user.UserId,
		SessionId:  user.SessionId,
		SessionKey: privateKeyString,
		PinToken:   user.PINTokenBase64,
		CreatedAt:  time.Now(),
	}

	for {
		err := key.setupPIN(ctx)
		if err == nil {
			break
		}
		log.Println(session.ServerError(ctx, err))
		time.Sleep(1 * time.Second)
	}
	for {
		err := key.setupOceanKey(ctx)
		if err == nil {
			break
		}
		log.Println(session.ServerError(ctx, err))
		time.Sleep(1 * time.Second)
	}

	for {
		err := key.persist(ctx)
		if err == nil {
			break
		}
		log.Println(session.TransactionError(ctx, err))
		time.Sleep(500 * time.Millisecond)
	}
	return key, nil
}

func (k *PoolKey) persist(ctx context.Context) error {
	stat := fmt.Sprintf("INSERT INTO pool_keys(%s) VALUES($1, $2, $3, $4, $5, $6, $7, $8)", strings.Join(poolKeysColumnsFull, ","))
	_, err := session.Database(ctx).Exec(ctx, stat, k.valuesFull()...)
	return err
}

func (k *PoolKey) setupOceanKey(ctx context.Context) error {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return session.ServerError(ctx, err)
	}
	oceanKey, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		return session.ServerError(ctx, err)
	}

	pub, err := x509.MarshalPKIXPublicKey(priv.Public())
	if err != nil {
		return session.ServerError(ctx, err)
	}
	sig := make([]byte, 140)
	handle := new(codec.MsgpackHandle)
	encoder := codec.NewEncoderBytes(&sig, handle)
	action := map[string][]byte{"U": pub}
	err = encoder.Encode(action)
	if err != nil {
		return session.ServerError(ctx, err)
	}

	k.SurfingKey = hex.EncodeToString(oceanKey)
	return nil
}

func generateSixDigitCode(ctx context.Context) (string, error) {
	var b [8]byte
	_, err := rand.Read(b[:])
	if err != nil {
		return "", err
	}
	c := binary.LittleEndian.Uint64(b[:]) % 1000000
	if c < 100000 {
		c = 100000 + c
	}
	return fmt.Sprint(c), nil
}

func (k *PoolKey) setupPIN(ctx context.Context) error {
	pin, err := generateSixDigitCode(ctx)
	if err != nil {
		return session.ServerError(ctx, err)
	}
	encryptedPIN, err := bot.EncryptEd25519PIN(ctx, pin, k.PinToken, k.SessionId, k.SessionKey, uint64(time.Now().UnixNano()))
	if err != nil {
		return err
	}
	err = bot.UpdatePin(ctx, "", encryptedPIN, k.UserId, k.SessionId, k.SessionKey)
	if err != nil {
		return err
	}
	encryptedPIN, encryptionHeader, err := encryptPIN(ctx, pin)
	if err != nil {
		return session.ServerError(ctx, err)
	}
	k.EncryptedPIN = encryptedPIN
	k.EncryptionHeader = encryptionHeader
	k.PlainPIN = pin
	return nil
}

func encryptPIN(ctx context.Context, pin string) (string, []byte, error) {
	aesKey := make([]byte, 32)
	_, err := rand.Read(aesKey)
	if err != nil {
		return "", nil, session.ServerError(ctx, err)
	}
	publicBytes, err := base64.StdEncoding.DecodeString(config.AssetPublicKey)
	if err != nil {
		return "", nil, session.ServerError(ctx, err)
	}
	assetPublicKey, err := x509.ParsePKCS1PublicKey(publicBytes)
	if err != nil {
		return "", nil, session.ServerError(ctx, err)
	}
	aesKeyEncrypted, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, assetPublicKey, aesKey, nil)
	if err != nil {
		return "", nil, session.ServerError(ctx, err)
	}
	encryptionHeader := make([]byte, encryptionHeaderLength)
	encryptionHeader = append(encryptionHeader, aesKeyEncrypted...)

	paddingSize := aes.BlockSize - len(pin)%aes.BlockSize
	paddingBytes := bytes.Repeat([]byte{byte(paddingSize)}, paddingSize)
	plainBytes := append([]byte(pin), paddingBytes...)
	cipherBytes := make([]byte, aes.BlockSize+len(plainBytes))
	iv := cipherBytes[:aes.BlockSize]
	_, err = rand.Read(iv)
	if err != nil {
		return "", nil, session.ServerError(ctx, err)
	}
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", nil, session.ServerError(ctx, err)
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherBytes[aes.BlockSize:], plainBytes)
	return base64.StdEncoding.EncodeToString(cipherBytes), encryptionHeader, nil
}
