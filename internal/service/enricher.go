package service

import (
	"encoding/json"
	"net/http"
)

type EnrichedData struct {
	Gender      string
	Age         int
	Nationality string
}

func EnrichPerson(name string) (EnrichedData, error) {
	var genderResp struct {
		Gender string `json:"gender"`
	}
	var ageResp struct {
		Age int `json:"age"`
	}
	var natResp struct {
		Country []struct {
			CountryID string `json:"country_id"`
		} `json:"country"`
	}

	getJSON := func(url string, target interface{}) error {
		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		return json.NewDecoder(resp.Body).Decode(target)
	}

	err := getJSON("https://api.genderize.io/?name="+name, &genderResp)
	if err != nil {
		return EnrichedData{}, err
	}
	err = getJSON("https://api.agify.io/?name="+name, &ageResp)
	if err != nil {
		return EnrichedData{}, err
	}
	err = getJSON("https://api.nationalize.io/?name="+name, &natResp)
	if err != nil {
		return EnrichedData{}, err
	}

	nationality := ""
	if len(natResp.Country) > 0 {
		nationality = natResp.Country[0].CountryID
	}

	return EnrichedData{
		Gender:      genderResp.Gender,
		Age:         ageResp.Age,
		Nationality: nationality,
	}, nil
}
