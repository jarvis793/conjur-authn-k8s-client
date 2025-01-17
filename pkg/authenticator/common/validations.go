package common

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/jarvis793/conjur-authn-k8s-client/pkg/log"
)

// ReadFileFunc defines the interface for reading an SSL Certificate from the env
type ReadFileFunc func(filename string) ([]byte, error)

func validTimeout(key, timeoutStr string) error {
	_, err := durationFromString(key, timeoutStr)
	return err
}

func validInt(key, value string) error {
	_, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf(log.CAKC060, key, value)
	}
	return nil
}

func validUsername(key, value string) error {
	if len(value) == 0 {
		return nil
	}
	_, err := NewUsername(value)
	return err
}

func ValidateSetting(key string, value string) error {
	switch key {
	case "CONJUR_AUTHN_LOGIN":
		return validUsername(key, value)
	case "CONJUR_CLIENT_CERT_RETRY_COUNT_LIMIT":
		return validInt(key, value)
	case "CONJUR_TOKEN_TIMEOUT":
		return validTimeout(key, value)
	case "JWT_TOKEN_PATH":
		return validatePath(value)
	default:
		return nil
	}
}

func ReadSSLCert(settings map[string]string, readFile ReadFileFunc) ([]byte, error) {
	SSLCert := settings["CONJUR_SSL_CERTIFICATE"]
	SSLCertPath := settings["CONJUR_CERT_FILE"]
	if SSLCert == "" && SSLCertPath == "" {
		return nil, errors.New(log.CAKC007)
	}

	if SSLCert != "" {
		return []byte(SSLCert), nil
	}
	return readFile(SSLCertPath)
}

func validatePath(path string) error {
	// Check if file already exists
	if _, err := os.Stat(path); err == nil {
		return nil
	}

	// Attempt to create the file and delete it right after
	var emptyData []byte
	if err := ioutil.WriteFile(path, emptyData, 0644); err == nil {
		os.Remove(path) // And delete it
		return nil
	}

	return fmt.Errorf(log.CAKC065, path)
}
