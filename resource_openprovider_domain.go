package main

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceDomains() *schema.Resource {
	return &schema.Resource{
		Create: CreateDomain,

		Schema: map[string]*schema.Schema{
			"domainname": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func CheckAvailability(domainnames string, m any) error {
	// Initializing client
	client := m.(Client)

	// Creating the request
	v := strings.Split(domainnames, ".")
	payload := map[string]any{
		"domains": []map[string]string{
			{
				"extension": v[1],
				"name":      v[0],
			},
		},
	}

	url := baseurl + "domains/check"

	req, err := client.NewRequest("POST", url, payload)
	if err != nil {
		return nil
	}

	res, err := client.Do(req)
	if err != nil {
		return nil
	}

	// If the returncode is anything other then 0, catch the error
	if returncode, ok := res["code"].(float64); ok {
		if returncode != 0 {
			// if desc, ok := res["desc"].(string); ok {
			// 	return nil
			// }
			return nil
		} else {
			return nil
		}
	} else {
		return nil
	}
}

func CreateDomain(d *schema.ResourceData, m any) error {
	// Initializing client
	client := m.(Client)

	// Check if domain is available
	domain_name := d.Get("domainname").(string)
	admin_handle := d.Get("admin_handle").(string)
	owner_handle := d.Get("owner_handle").(string)
	tech_handle := d.Get("tech_handle").(string)
	period := d.Get("period").(string)
	name_servers := d.Get("name_servers").(string)
	autorenew := d.Get("autorenew").(string)

	rvalue := CheckAvailability(domain_name, m)
	if rvalue != nil {
		return rvalue
	}

	domain_name_split := strings.Split(domain_name, ".")

	payload := map[string]any{
		"admin_handle": admin_handle,
		"owner_handle": owner_handle,
		"tech_handle":  tech_handle,
		"domain": map[string]string{
			"name":      domain_name_split[0],
			"extension": domain_name_split[1],
		},
		"period": period,
		"nameservers": []map[string]any{
			{"name": name_servers},
		},
		"autorenew": autorenew,
	}

	url := baseurl + "domains"

	req, err := client.NewRequest("POST", url, payload)
	if err != nil {
		return nil
	}

	res, err := client.Do(req)
	if err != nil {
		return nil
	}

	// If the returncode is anything other then 0, catch the error
	if returncode, ok := res["code"].(float64); ok {
		if returncode != 0 {
			// if desc, ok := res["desc"].(string); ok {
			// 	return nil
			// }
			return nil
		} else {
			return nil
		}
	} else {
		return nil
	}
}
