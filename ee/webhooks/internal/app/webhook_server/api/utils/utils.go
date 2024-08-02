package utils

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type ErrorType string

const (
	NoneType       ErrorType = "NONE"
	ValidationType ErrorType = "VALIDATION"
	NotFoundType   ErrorType = "NOT_FOUND"
	InternalType   ErrorType = "INTERNAL"
)

type Response[T interface{}] struct {
	T    ErrorType
	Err  error
	Data *T
}

func SuccessResp[T interface{}](d T) Response[T] {
	return Response[T]{
		T:    NoneType,
		Err:  nil,
		Data: &d,
	}
}

func ValidationErrorResp[T interface{}](err error) Response[T] {
	return Response[T]{
		T:    ValidationType,
		Err:  err,
		Data: nil,
	}
}

func InternalErrorResp[T interface{}](err error) Response[T] {
	return Response[T]{
		T:    InternalType,
		Err:  err,
		Data: nil,
	}
}

func NotFoundErrorResp[T interface{}](err error) Response[T] {
	return Response[T]{
		T:    NotFoundType,
		Err:  err,
		Data: nil,
	}
}

const (
	ErrValidation  = "VALIDATION_TYPE"
	ErrHealthcheck = "HEALTHCHECK_STATUS"
)

func DecodeJSONBody(r *http.Request, v any) error {
	return json.NewDecoder(r.Body).Decode(&v)

}

var (
	ErrInvalidEndpoint   = errors.New("endpoint should be a valid url")
	ErrInvalidEventTypes = errors.New("eventTypes should be filled")
	ErrInvalidSecret     = errors.New("decoded secret should be of size 24")
)

func ValidateEndpoint(endpoint string) error {
	if u, err := url.Parse(endpoint); err != nil || len(u.String()) == 0 {
		return ErrInvalidEndpoint
	}
	return nil
}

func ValidateSecret(secret *string) error {

	if *secret != "" {
		var decoded []byte
		var err error
		if decoded, err = base64.StdEncoding.DecodeString(*secret); err != nil {
			return fmt.Errorf("secret should be base64 encoded: %w", err)
		}

		if len(decoded) != 24 {
			return fmt.Errorf("decoded secret should have 24 caracters")
		}
	} else {
		*secret = newSecret()
	}
	return nil
}

func FormatEvents(events *[]string) error {

	for i, t := range *events {
		if len(t) == 0 {
			return ErrInvalidEventTypes
		}
		(*events)[i] = strings.ToLower(t)
	}

	return nil
}

func newSecret() string {
	token := make([]byte, 24)
	_, err := rand.Read(token)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(token)
}

func NewSecret() string {
	return newSecret()
}

func ReadCursor(strCursor string) (int, error) {
	if strCursor == "" {
		return 0, nil
	}
	return strconv.Atoi(strCursor)
}

func PaginationCursor(cursor int, hasMore bool) (previous, next string) {
	strP := " "
	strN := " "

	if hasMore {
		strN = fmt.Sprintf("%d", cursor+1)
	}
	if cursor > 0 {
		strP = fmt.Sprintf("%d", cursor-1)
	}

	return strP, strN
}

type Secret struct {
	Secret string `json:"secret"`
}

type Endpoint struct {
	Endpoint string `json:"endpoint"`
}

type Retry struct {
	Retry bool `json:"retry"`
}

func ToValues[T []*G, G any](in T) []G {
	out := make([]G, 0)
	for _, v := range in {
		out = append(out, *v)
	}
	return out
}

func Pagination(page int, pageSize int) (startPage int, endPage int) {

	startPage = page * pageSize
	endPage = ((page + 1) * pageSize) - 1

	return startPage, endPage
}
