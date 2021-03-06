package provider

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go/service/lexmodelsv2"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{},
		ResourcesMap: map[string]*schema.Resource{
			"awslexbotv2_bot":       resourceBot(),
			"awslexbotv2_uploadurl": resourceUploadUrl(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			// "awslexbotv2_instance": dataSourceBot(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

// Client -
type Client struct {
	STSClient      *sts.STS
	LexBotV2Client *lexmodelsv2.LexModelsV2
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	log.Println("[DEBUG] lexbotv2 provider configuree")
	c := Client{
		LexBotV2Client: lexModelsV2Service(),
		STSClient:      stsService(),
	}

	return c, diags
}
