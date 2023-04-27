package utils

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/middleware"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
	"hash"
	"math"
	"math/rand"
	"net/http"
	"strings"
	"sync/atomic"
)

var (
	RequiredFieldMessage = func(fields ...string) string {
		return fmt.Sprintf("%v required", fields)
	}
	reqid uint64
)

const (
	RealUserIpKey    = "X-Original-Forwarded-For"
	AdminUserKey     = "admin"
	RoleKey          = "role"
	KeyForSecretKey  = "Secret-Key"
	AuthorizationKey = "Authorization"
	SaltSize         = 16
	myCharset        = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ" + "0123456789"
	float64EqualityThreshold = 1e-9
)

func GetTracingID(ctx context.Context) string {
	return middleware.GetReqID(ctx)
}

func GetHasher() hash.Hash {
	return md5.New()
}

func GenerateSalt() *string {
	var salt = make([]rune, SaltSize)
	for i := range salt {
		salt[i] = rune(myCharset[rand.Intn(len(myCharset))])
	}
	saltStr := string(salt)
	return &saltStr
}

func HashPassword(password string, salt string) string {
	passwordBytes := []byte(password)
	md5Hasher := GetHasher()

	passwordBytes = append(passwordBytes, salt...)
	md5Hasher.Write(passwordBytes)

	hashedPassInBytes := md5Hasher.Sum(nil)
	hashedPassInHex := hex.EncodeToString(hashedPassInBytes)

	return hashedPassInHex
}

func IsSamePassword(hashedPassword, currentPassword string, salt string) bool {
	currentPasswordHash := HashPassword(currentPassword, salt)
	return hashedPassword == currentPasswordHash
}

func GetKeyFromHeader(r *http.Request, key string) string {
	return r.Header.Get(key)
}

func GetAuthTokenFromHeader(r *http.Request) string {
	bearer := GetKeyFromHeader(r, AuthorizationKey)
	return strings.TrimPrefix(bearer, "Bearer ")
}

/*-----------------------------------------------*/

func BoolP(boolValue bool) *bool {
	return &boolValue
}

func CustomJsonMarshal(data interface{}, tag string) ([]byte, error) {
	var json = jsoniter.Config{
		EscapeHTML:             true,
		SortMapKeys:            true,
		ValidateJsonRawMessage: true,
		TagKey:                 tag,
	}.Froze()

	return json.Marshal(data)
}

func GetSlug(str string) string {
	str = strings.TrimSpace(strings.ToLower(str))
	return strings.Join(strings.Split(str, " "), "-")
}

func SetTracingID(ctx context.Context) context.Context {
	uid := uuid.New().String()
	myid := atomic.AddUint64(&reqid, 1)
	requestID := fmt.Sprintf("%s-%06d", uid, myid)
	ctx = context.WithValue(ctx, middleware.RequestIDKey, requestID)
	return ctx
}

func DecodeInterface(input, output interface{}) error {
	data, err := json.Marshal(input)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, output)
}

func AlmostEqual(a, b float64) bool {
	return math.Abs(a-b) <= float64EqualityThreshold
}
