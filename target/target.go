package target

import (
	"errors"
	"io/ioutil"
	"net/http"
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
	response, err := http.Get(target)

	if err != nil {
		return false, err
	}
	defer response.Body.Close()
	if response.StatusCode/100 != 2 {
		return false, errors.New("Status code returned was " + response.Status)
	}
	var body []byte
	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return false, err
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
