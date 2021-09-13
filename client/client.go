package client

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

const (
	RegisterServer = "https://version-check.nuclei.sh/"
	VersionsCall   = "versions"
	IgnoreCall     = "ignore"
)

var nucleiVersion string

// LatestVersion is the latest version info for nuclei and templates repos
type LatestVersion struct {
	Nuclei    string
	Templates string
}

func InitNucleiVersion(version string) {
	nucleiVersion = fmt.Sprintf("Nuclei/%s", version)
}

// GetLatestNucleiTemplatesVersion returns the latest version info for nuclei and templates repos
func GetLatestNucleiTemplatesVersion() (*LatestVersion, error) {
	body, err := callRegisterServer(VersionsCall)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	data := make(map[string]string)
	if err := jsoniter.NewDecoder(body).Decode(&data); err != nil {
		return nil, err
	}
	return &LatestVersion{Nuclei: data["nuclei"], Templates: data["templates"]}, nil
}

// GetLatestIgnoreFile returns the latest version of nuclei ignore
func GetLatestIgnoreFile() ([]byte, error) {
	body, err := callRegisterServer(IgnoreCall)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	data, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// callRegisterServer makes a request to RegisterServer with a call.
func callRegisterServer(call string) (io.ReadCloser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, RegisterServer+call, nil)
	if err != nil {
		return nil, errors.Wrap(err, "could not make request")
	}
	if nucleiVersion != "" {
		req.Header.Set("User-Agent", nucleiVersion)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "could not do request")
	}
	return resp.Body, nil
}
