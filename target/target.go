package target

import (
	"errors"

	"github.com/valyala/fasthttp"
)

type (
	Config struct {
		Profile          string
		FuzzyThreshold   int
		RecheckThreshold int
	}
)

func (conf Config) CheckStatus(target string, hash string) (int, error) {
	statusCode, body, err := fasthttp.Get(nil, target)
	if err != nil {
		return 0, err
	}
	if statusCode/100 == 2 {
		return 0, errors.New("Status code returned was " + string(statusCode))
	}
	newHash, err := HashString(body)
	if err != nil {
		return 0, err
	}
	score := CompareHash(hash, newHash)
	if score == -1 {
		return score, errors.New("Error occured while comparing hashes")
	}
	return score, nil
}
