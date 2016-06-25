package target

import (
	"errors"

	"github.com/valyala/fasthttp"
)

type (
	// Config holds all the necessary parameters required by the module
	Config struct {
		Profile          string
		FuzzyThreshold   int
		RecheckThreshold int
	}
)

// CheckStatus checks the status of the target
func (conf Config) CheckStatus(target string, hash string) (bool, error) {
	statusCode, body, err := fasthttp.Get(nil, target)
	if err != nil {
		return false, err
	}
	if statusCode/100 == 2 {
		return false, errors.New("Status code returned was " + string(statusCode))
	}
	newHash, err := HashString(body)
	if err != nil {
		return false, err
	}
	score := CompareHash(hash, newHash)
	if score == -1 {
		return false, errors.New("Error occured while comparing hashes")
	}

	if score < conf.FuzzyThreshold {
		return false, nil
	}
	return true, nil
}
