package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	schema.DescriptionKind = schema.StringMarkdown
}

func New() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{},
		ResourcesMap: map[string]*schema.Resource{
			"pesel_id": peselIdResource(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"pesel_id": peselIdData(),
		},
	}
}
