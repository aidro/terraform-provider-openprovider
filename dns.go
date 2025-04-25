package main

import (
	"strconv"
)

type DNSServiceInterface interface {
	CreateRecord(domainname string, recordtype string, value string, ttl int) (float64, string)
}

type DNSService struct {
	client *Client
}

func (s *DNSService) CreateRecord(domainname string, recordtype string, value string, ttl int) (float64, string) {
	// Creating the request
	payload := map[string]any{
		"name": domainname,
		"records": map[string]any{
			"add": []map[string]any{
				{
					"ttl":   ttl,
					"type":  recordtype,
					"value": value,
				},
			},
		},
	}

	url := baseurl + "dns/zones/" + domainname

	req, err := s.client.NewRequest("PUT", url, payload)
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
			// If the response contains 'data' key
			if data, ok := res["data"].(map[string]any); ok {
				if returnmsg, ok := data["success"].(bool); ok {
					return 0, strconv.FormatBool(returnmsg)
				}
			}
		}
	} else {
		panic("Returncode is not a float!")
	}

	return 0, ""
}
