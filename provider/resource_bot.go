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
	connectSvc := meta.(Client).LexBotV2Client

	params := &lexmodelsv2.CreateBotInput{
		BotName:                 aws.String(d.Get("name").(string)),
		Description:             aws.String(d.Get("description").(string)),
		IdleSessionTTLInSeconds: aws.Int64(d.Get("idle_session_ttl").(int64)),
	}

	resp, err := connectSvc.CreateBot(params)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(aws.StringValue(resp.BotId))
	// d.Set("bot_id", aws.StringValue(resp.BotId))
	// d.Set("arn", aws.StringValue(resp.Arn))

	time.Sleep(3 * time.Minute) // wait 3m for connect instance creation

	return diags
}

func resourceBotRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	// connectSvc := meta.(Client).LexBotV2Client

	// instanceID := d.Get("instance_id").(string)

	// params := &connect.DescribeBotInput{
	// 	BotId: aws.String(instanceID),
	// }
	// resp, err := connectSvc.DescribeBot(params)
	// if err != nil {
	// 	return diag.FromErr(err)
	// }

	// d.SetId(instanceID)
	// d.Set("arn", aws.StringValue(resp.Bot.Arn))
	// d.Set("instance_alias", aws.StringValue(resp.Bot.BotAlias))
	// d.Set("identity_management_type", aws.StringValue(resp.Bot.IdentityManagementType))
	// d.Set("inbound_calls_enabled", aws.BoolValue(resp.Bot.InboundCallsEnabled))
	// d.Set("outbound_calls_enabled", aws.BoolValue(resp.Bot.OutboundCallsEnabled))

	return diags
}

func resourceBotUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	// connectSvc := meta.(Client).LexBotV2Client

	// instanceID := aws.String(d.Id())

	// if d.HasChange("inbound_calls_enabled") {
	// 	params := &connect.UpdateBotAttributeInput{
	// 		BotId:         instanceID,
	// 		AttributeType: aws.String("INBOUND_CALLS"),
	// 		Value:         aws.String(strconv.FormatBool(d.Get("inbound_calls_enabled").(bool))),
	// 	}
	// 	_, err := connectSvc.UpdateBotAttribute(params)
	// 	if err != nil {
	// 		return diag.FromErr(err)
	// 	}
	// }
	// if d.HasChange("outbound_calls_enabled") {
	// 	params := &connect.UpdateBotAttributeInput{
	// 		BotId:         instanceID,
	// 		AttributeType: aws.String("OUTBOUND_CALLS"),
	// 		Value:         aws.String(strconv.FormatBool(d.Get("outbound_calls_enabled").(bool))),
	// 	}
	// 	_, err := connectSvc.UpdateBotAttribute(params)
	// 	if err != nil {
	// 		return diag.FromErr(err)
	// 	}
	// }

	return diags
}

func resourceBotDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	// connectSvc := meta.(Client).LexBotV2Client

	// params := &connect.DeleteBotInput{
	// 	BotId: aws.String(d.Id()),
	// }

	// _, err := connectSvc.DeleteBot(params)
	// if err != nil {
	// 	return diag.FromErr(err)
	// }

	// // d.SetId("") is automatically called assuming delete returns no errors, but
	// // it is added here for explicitness.
	// d.SetId("")

	return diags
}
