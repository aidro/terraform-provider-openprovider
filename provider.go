package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "username of account",
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "password of account",
			},
		},

		ConfigureFunc: providerConfigure,

		ResourcesMap: map[string]*schema.Resource{
			"openprovider_domain": resourceDomains(),
		},
	}
}

func providerConfigure(d *schema.ResourceData) (any, error) {
	// load username and password
	username := d.Get("username").(string)
	password := d.Get("password").(string)

	if username != "" && password != "" {
		client, err := NewClient(username, password)
		if err != nil {
			return client, nil
		}
	}

	return nil, nil
}
