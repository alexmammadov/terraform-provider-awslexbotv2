package provider

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lexmodelsv2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceBot() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBotCreate,
		ReadContext:   resourceBotRead,
		UpdateContext: resourceBotUpdate,
		DeleteContext: resourceBotDelete,

		Schema: map[string]*schema.Schema{
			"name":             {Type: schema.TypeString, Required: true},
			"description":      {Type: schema.TypeString, Required: true},
			"idle_session_ttl": {Type: schema.TypeInt, Required: true},
			"id":               {Type: schema.TypeString, Computed: true, Optional: true},
		},
	}
}

func resourceBotCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	svc := meta.(Client).LexBotV2Client

	params := &lexmodelsv2.CreateBotInput{
		BotName:                 aws.String(d.Get("name").(string)),
		Description:             aws.String(d.Get("description").(string)),
		IdleSessionTTLInSeconds: aws.Int64(int64(d.Get("idle_session_ttl").(int))),
		DataPrivacy: &lexmodelsv2.DataPrivacy{
			ChildDirected: aws.Bool(false),
		},
		RoleArn: aws.String("arn:aws:iam::021721647551:role/service-role/test-lex-role-2te64zgs"),
	}

	resp, err := svc.CreateBot(params)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(aws.StringValue(resp.BotId))
	// d.Set("bot_id", aws.StringValue(resp.BotId))
	// d.Set("arn", aws.StringValue(resp.Arn))

	time.Sleep(10 * time.Second)

	return diags
}

func resourceBotRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	svc := meta.(Client).LexBotV2Client

	id := d.Id()

	params := &lexmodelsv2.DescribeBotInput{
		BotId: aws.String(id),
	}
	resp, err := svc.DescribeBot(params)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)
	d.Set("name", aws.StringValue(resp.BotName))
	d.Set("description", aws.StringValue(resp.Description))
	d.Set("idle_session_ttl", int(aws.Int64Value(resp.IdleSessionTTLInSeconds)))

	return diags
}

func resourceBotUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	svc := meta.(Client).LexBotV2Client

	id := aws.String(d.Id())

	params := &lexmodelsv2.UpdateBotInput{
		BotId:                   id,
		BotName:                 aws.String(d.Get("name").(string)),
		Description:             aws.String(d.Get("description").(string)),
		IdleSessionTTLInSeconds: aws.Int64(int64(d.Get("idle_session_ttl").(int))),
		DataPrivacy: &lexmodelsv2.DataPrivacy{
			ChildDirected: aws.Bool(false),
		},
		RoleArn: aws.String("arn:aws:iam::021721647551:role/service-role/test-lex-role-2te64zgs"),
	}
	_, err := svc.UpdateBot(params)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceBotDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	svc := meta.(Client).LexBotV2Client

	params := &lexmodelsv2.DeleteBotInput{
		BotId:                  aws.String(d.Id()),
		SkipResourceInUseCheck: aws.Bool(true),
	}

	_, err := svc.DeleteBot(params)
	if err != nil {
		return diag.FromErr(err)
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	time.Sleep(10 * time.Second)

	return diags
}
