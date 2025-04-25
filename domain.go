package main

import (
	"strings"
)

type DomainServiceInterface interface {
	CheckAvailability(domainnames []string) (float64, string)
	CreateDomain(admin_handle string, owner_handle string, tech_handle string, domainname []string, period int, name_servers []string, autorenew string) (float64, string)
}

type DomainService struct {
	client *Client
}

func (s *DomainService) CheckAvailability(domainnames []string) (float64, string) {
	// Creating the request
	v := strings.Split(domainnames[0], ".")
	payload := map[string]any{
		"domains": []map[string]string{
			{
				"extension": v[1],
				"name":      v[0],
			},
		},
	}

	url := baseurl + "domains/check"

	req, err := s.client.NewRequest("POST", url, payload)
	if err != nil {
		panic(err)
	}

	res, err := s.client.Do(req)
	if err != nil {
		panic(err)
	}

	// If the returncode is anything other then 0, catch the error
	if returncode, ok := res["code"].(float64); ok {
		if returncode != 0 {
			if desc, ok := res["desc"].(string); ok {
				return returncode, desc
			}
		} else {
			return returncode, "nil"
		}
	} else {
		return returncode, "nil"
	}

	return 1, "nil"
}

func (s *DomainService) CreateDomain(admin_handle string, owner_handle string, tech_handle string, domainname []string, period int, name_servers []string, autorenew string) (float64, string) {
	// Check if domain is available
	rvalue, errortje := s.CheckAvailability(domainname)
	if rvalue != 0 {
		// If domainname is not free, return error
		return rvalue, errortje
	}
	v := strings.Split(domainname[0], ".")

	payload := map[string]any{
		"admin_handle": admin_handle,
		"owner_handle": owner_handle,
		"tech_handle":  tech_handle,
		"domain": map[string]string{
			"name":      v[0],
			"extension": v[1],
		},
		"period": period,
		"nameservers": []map[string]any{
			{"name": name_servers},
		},
		"autorenew": autorenew,
	}

	url := baseurl + "domains"

	req, err := s.client.NewRequest("POST", url, payload)
	if err != nil {
		panic(err)
	}

	res, err := s.client.Do(req)
	if err != nil {
		panic(err)
	}

	// If the returncode is anything other then 0, catch the error
	if returncode, ok := res["code"].(float64); ok {
		if returncode != 0 {
			if desc, ok := res["desc"].(string); ok {
				return returncode, desc
			}
		} else {
			return returncode, "nil"
		}
	} else {
		return returncode, "nil"
	}

	return 1, "nil"
}
