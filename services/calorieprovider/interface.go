package calorieprovider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type Provider interface {
	GetCalories(meal string) (int, error)
}

type ProviderImpl struct {
}

const (
	appID  = "d7684a71"
	appKey = "42d601cad14f572657f5d1d3293e2fdf"
)

type queryPayload struct {
	Query string `json:"query"`
}

func (p *ProviderImpl) GetCalories(meal string) (int, error) {
	url := "https://trackapi.nutritionix.com/v2/natural/nutrients"
	payload, err := json.Marshal(&queryPayload{Query: meal})
	if err != nil {
		return -1, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return -1, err
	}

	req.Header.Set("x-app-id", appID)
	req.Header.Set("x-app-key", appKey)
	req.Header.Set("x-remote-user-id", "0") // set to 0 for development mode

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()

	bufferSize := int64(1024 * 1024)
	responseBytes, err := ioutil.ReadAll(io.LimitReader(resp.Body, bufferSize))
	if err != nil {
		return -1, err
	}

	response := map[string][]map[string]interface{}{}
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		return -1, err
	}

	calories, ok := response["foods"][0]["nf_calories"].(float64)
	if !ok {
		return -1, fmt.Errorf("Failed to parse calories from provider")
	}
	return int(calories), nil
}
