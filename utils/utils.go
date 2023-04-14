package utils

import (
	"context"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/middleware"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/viper"
	"go-rest-api/model"
	"hash"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strings"
	"sync/atomic"
	"time"
)

const (
	BetaClusterRedisURL = "beta-redis-master.redis.svc.cluster.local:6379"
	SuccessMessage      = "successful"
	DeletedSuccessfully = "deleted successfully"
	CreatedSuccessfully = "created successfully"
)

var (
	RequiredFieldMessage = func(fields ...string) string {
		return fmt.Sprintf("%v required", fields)
	}
	reqid uint64
)

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

func GetTracingID(ctx context.Context) string {
	return middleware.GetReqID(ctx)
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

const myCharset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

//var seededRand *rand.Rand = rand.New(
//	rand.NewSource(time.Now().UnixNano()))
//
//func RandStr(length int) string {
//
//	b := make([]byte, length)
//	for i := range b {
//		b[i] = myCharset[seededRand.Intn(len(myCharset))]
//	}
//	return string(b)
//}

func DecodeInterface(input, output interface{}) error {
	data, err := json.Marshal(input)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, output)
}

const float64EqualityThreshold = 1e-9

func AlmostEqual(a, b float64) bool {
	return math.Abs(a-b) <= float64EqualityThreshold
}

/*-----------------------------------------------*/

const (
	RealUserIpKey    = "X-Original-Forwarded-For"
	AdminUserKey     = "admin"
	RoleKey          = "role"
	KeyForSecretKey  = "Secret-Key"
	AuthorizationKey = "Authorization"
	SaltSize         = 16
)

func GetHasher() hash.Hash {
	return md5.New()
}

func GenerateSalt() *string {
	var salt = make([]byte, 16)
	_, err := rand.Read(salt)

	if err != nil {
		log.Println("error generating salt: ", err.Error())
		return nil
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

func GetUserByJwtToken(jwtTkn string) (*model.AuthUserData, error) {
	sendError := errors.New("failed to validate token")

	client := http.Client{Timeout: time.Minute * 2}
	viper.AutomaticEnv()
	authBaseUrl := viper.GetString("app.auth_base_url")

	authUrl := authBaseUrl + "/api/v1/auth/validate-token"

	req, err := http.NewRequest(http.MethodGet, authUrl, nil)
	if err != nil {
		log.Println(err)
		return nil, sendError
	}

	req.Header.Set("Authorization", jwtTkn)

	log.Println("sending query:")
	// send the request
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return nil, sendError
	}
	log.Println("query sent: reading resp: status:", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		log.Println(resp.StatusCode)
		fmt.Println(resp.Status)
		return nil, sendError
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	err = resp.Body.Close()
	if err != nil {
		log.Println(err)
		return nil, sendError
	}
	var reply model.AuthUserResponse

	err = json.Unmarshal(body, &reply)

	res := reply.Data
	if res.UserName == "" {
		return nil, sendError
	}

	return res, nil
}
