package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"testing"
)

var Providers map[string]func() (*schema.Provider, error)
var Provider *schema.Provider

func init() {
	Provider = New()
	Providers = map[string]func() (*schema.Provider, error){
		"pesel": func() (*schema.Provider, error) {
			return Provider, nil
		},
	}
}

func PreCheck(t *testing.T) {}
