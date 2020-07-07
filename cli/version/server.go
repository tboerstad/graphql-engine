package version

import (
	"io/ioutil"
	"net/http"
	"os"

	yaml "github.com/ghodss/yaml"
	"github.com/pkg/errors"
)

type serverVersionResponse struct {
	Version string `json:"version"`
}

// FetchServerVersion reads the version from server.
func FetchServerVersion(endpoint string, client *http.Client) (version string, err error) {
	req, err := http.NewRequest("GET", endpoint, nil)
	
	token, success := os.LookupEnv("AUTH_TOKEN")
	if success {
		req.Header.Set("Authorization", "Bearer " + token)
	}	

	response, err := client.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "failed making version api call")
	}
	if response.StatusCode != http.StatusOK {
		switch response.StatusCode {
		case http.StatusNotFound:
			return "", nil
		default:
			return "", errors.Errorf("GET %s failed - [%d]", endpoint, response.StatusCode)
		}
	} else {
		defer response.Body.Close()
		data, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return "", errors.Wrap(err, "cannot read version api response")
		}
		var v serverVersionResponse
		err = yaml.Unmarshal(data, &v)
		if err != nil {
			return "", errors.Wrap(err, "failed to parse version api response")
		}
		return v.Version, nil
	}
}
