package provider

import (
	"context"

	"github.com/aws/aws-sdk-go/service/lexmodelsv2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{},
		ResourcesMap: map[string]*schema.Resource{
			"awslexbotv2_bot": resourceBot(),
			// "awslexbotv2_instance_lex_bot":      resourceBotLexBot(),
			// "awslexbotv2_instance_contact_flow": resourceBotContactFlow(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			// "awslexbotv2_instance": dataSourceBot(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

// Client -
type Client struct {
	LexBotV2Client *lexmodelsv2.LexModelsV2
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	c := Client{
		LexBotV2Client: lexModelsV2Service(),
	}

	return c, diags
}
