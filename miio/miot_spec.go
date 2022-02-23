package miio

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const (
	urlAll string = "http://miot-spec.org/miot-spec-v2/instances?status=all"
	// get parameter should be taken from field 'type' Instance structure
	urlDetails string = "http://miot-spec.org/miot-spec-v2/instance?type="
)

type Instance struct {
	Status  string `json:"status"`
	Model   string `json:"model"`
	Version int    `json:"version"`
	Type    string `json:"type"`
}

type Instances struct {
	Instances []Instance `json:"instances"`
}

func GetInstances() (*Instances, error) {
	response, err := http.Get(urlAll)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	jsonResponse, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var i *Instances
	if err := json.Unmarshal(jsonResponse, &i); err != nil {
		return nil, err
	}

	return i, nil
}

func (i *Instances) String() string {
	b, err := json.Marshal(i)
	if err != nil {
		return ""
	}
	return string(b)
}

type Details struct {
	Type     string    `json:"type"`
	Desc     string    `json:"description"`
	Services []Service `json:"services"`
}

func GetDetail(urn string) (*Details, error) {
	var url string = urlDetails + urn
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	res, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var detail *Details
	if err := json.Unmarshal(res, &detail); err != nil {
		return nil, err
	}

	return detail, nil
}

func (a *Details) String() string {
	b, err := json.Marshal(a)
	if err != nil {
		return ""
	}
	return string(b)
}

type Service struct {
	Id    int        `json:"id"`
	Type  string     `json:"type"`
	Desc  string     `json:"description"`
	Props []Property `json:"properties"`
}

func (a *Service) String() string {
	b, err := json.Marshal(a)
	if err != nil {
		return ""
	}
	return string(b)
}

type Property struct {
	Id     int      `json:"id"`
	Type   string   `json:"type"`
	Desc   string   `json:"description"`
	Format string   `json:"format"`
	Access []string `json:"access"`
}

func (a *Property) String() string {
	b, err := json.Marshal(a)
	if err != nil {
		return ""
	}
	return string(b)
}
