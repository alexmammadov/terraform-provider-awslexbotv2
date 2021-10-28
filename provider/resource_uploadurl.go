package provider

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lexmodelsv2"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceUploadUrl() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUploadUrlCreate,
		ReadContext:   resourceUploadUrlRead,
		UpdateContext: resourceUploadUrlUpdate,
		DeleteContext: resourceUploadUrlDelete,

		Schema: map[string]*schema.Schema{
			"file_path":        {Type: schema.TypeString, Required: true},
			"bot_name":         {Type: schema.TypeString, Required: true},
			"idle_session_ttl": {Type: schema.TypeInt, Required: true},
			"role_arn":         {Type: schema.TypeString, Required: true},
			"etag":             {Type: schema.TypeString, Required: true},
			"lambda_arn":       {Type: schema.TypeString, Optional: true},
			"import_id":        {Type: schema.TypeString, Computed: true, Optional: true},
			"upload_url":       {Type: schema.TypeString, Computed: true, Optional: true},
			"bot_alias_arn":    {Type: schema.TypeString, Computed: true, Optional: true},
		},
	}
}

func resourceUploadUrlCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	svc := meta.(Client).LexBotV2Client
	log.Println("[DEBUG] calling createuploadurl")

	urlResp, err := svc.CreateUploadUrl(&lexmodelsv2.CreateUploadUrlInput{})
	if err != nil {
		log.Println("[DEBUG] calling createuploadurl error")

		return diag.FromErr(err)
	}
	log.Println("[DEBUG] calling createuploadurl success")

	d.SetId(aws.StringValue(urlResp.ImportId))
	d.Set("import_id", aws.StringValue(urlResp.ImportId))
	d.Set("upload_url", aws.StringValue(urlResp.UploadUrl))

	// time.Sleep(10 * time.Second)
	log.Println("[DEBUG] calling uploadfile")

	err = uploadFile(aws.StringValue(urlResp.UploadUrl), d.Get("file_path").(string))
	if err != nil {
		log.Println("[DEBUG] calling uploadfile error")

		return diag.FromErr(err)
	}
	log.Println("[DEBUG] calling uploadfile success")
	log.Println("[DEBUG] calling start import")

	_, err = svc.StartImport(&lexmodelsv2.StartImportInput{
		ImportId:      urlResp.ImportId,
		MergeStrategy: aws.String(lexmodelsv2.MergeStrategyOverwrite),
		ResourceSpecification: &lexmodelsv2.ImportResourceSpecification{
			BotImportSpecification: &lexmodelsv2.BotImportSpecification{
				BotName:                 aws.String(d.Get("bot_name").(string)),
				DataPrivacy:             &lexmodelsv2.DataPrivacy{ChildDirected: aws.Bool(false)},
				RoleArn:                 aws.String(d.Get("role_arn").(string)),
				IdleSessionTTLInSeconds: aws.Int64(int64(d.Get("idle_session_ttl").(int))),
			},
		},
	})
	if err != nil {
		log.Println("[DEBUG] calling start import error")

		return diag.FromErr(err)
	}
	log.Println("[DEBUG] calling start import success")
	log.Println("[DEBUG] calling describe import")

	respDescr, err := svc.DescribeImport(&lexmodelsv2.DescribeImportInput{
		ImportId: urlResp.ImportId,
	})
	if err != nil {
		log.Println("[DEBUG] calling describe import error")

		return diag.FromErr(err)
	}
	log.Println("[DEBUG] calling describe success")

	botId := respDescr.ImportedResourceId
	log.Println("[DEBUG] calling ListBotAliases")
	time.Sleep(10 * time.Second) // wait 10s for completeleting current task

	respAliases, err := svc.ListBotAliases(&lexmodelsv2.ListBotAliasesInput{
		BotId: botId,
	})
	if err != nil {
		log.Println("[DEBUG] calling ListBotAliases error")

		return diag.FromErr(err)
	}
	log.Println("[DEBUG] calling ListBotAliases success")

	if len(respAliases.BotAliasSummaries) == 0 {
		return diag.FromErr(errors.New("no aliases for the bot"))
	}
	log.Println("[DEBUG] calling ListBotAliases count success")

	aliasId := respAliases.BotAliasSummaries[0].BotAliasId
	region := svc.Config.Region

	stsSvc := meta.(Client).STSClient
	log.Println("[DEBUG] calling GetCallerIdentityRequest")

	reqAcc, respAcc := stsSvc.GetCallerIdentityRequest(&sts.GetCallerIdentityInput{})
	err = reqAcc.Send()
	if err != nil {
		log.Println("[DEBUG] calling GetCallerIdentityRequest error")

		return diag.FromErr(err)
	}
	log.Println("[DEBUG] calling GetCallerIdentityRequest success")

	aliasArn := fmt.Sprintf("arn:aws:lex:%s:%s:bot-alias/%s/%s",
		aws.StringValue(region),
		aws.StringValue(respAcc.Account),
		aws.StringValue(botId),
		aws.StringValue(aliasId))

	d.Set("bot_alias_arn", aliasArn)

	lambdaArn := d.Get("lambda_arn").(string)
	lambdaVersion := "1.0"

	if len(lambdaArn) > 0 {
		time.Sleep(30 * time.Second)
		reqAlias, respAlias := svc.DescribeBotAliasRequest(&lexmodelsv2.DescribeBotAliasInput{
			BotId:      botId,
			BotAliasId: aliasId,
		})
		err = reqAlias.Send()
		if err != nil {
			log.Println("[DEBUG] calling DescribeBotAliasRequest error")
			return diag.FromErr(err)
		}
		log.Println("[DEBUG] calling DescribeBotAliasRequest success")

		respAlias.BotAliasLocaleSettings = make(map[string]*lexmodelsv2.BotAliasLocaleSettings)
		respAlias.BotAliasLocaleSettings["en_US"] = &lexmodelsv2.BotAliasLocaleSettings{
			Enabled: aws.Bool(true),
			CodeHookSpecification: &lexmodelsv2.CodeHookSpecification{
				LambdaCodeHook: &lexmodelsv2.LambdaCodeHook{
					LambdaARN:                aws.String(lambdaArn),
					CodeHookInterfaceVersion: aws.String(lambdaVersion),
				},
			},
		}
		reqUpdateAlias, respUpdateAlias := svc.UpdateBotAliasRequest((&lexmodelsv2.UpdateBotAliasInput{
			BotAliasId:                respAlias.BotAliasId,
			BotAliasLocaleSettings:    respAlias.BotAliasLocaleSettings,
			BotAliasName:              respAlias.BotAliasName,
			BotId:                     respAlias.BotId,
			BotVersion:                respAlias.BotVersion,
			ConversationLogSettings:   respAlias.ConversationLogSettings,
			Description:               respAlias.Description,
			SentimentAnalysisSettings: respAlias.SentimentAnalysisSettings,
		}))
		err = reqUpdateAlias.Send()
		if err != nil {
			log.Println("[DEBUG] calling UpdateBotAliasRequest error")
			return diag.FromErr(err)
		}
		log.Println("[DEBUG] calling UpdateBotAliasRequest success", respUpdateAlias)
	}

	return diags
}

func uploadFile(serverURL, filename string) error {
	// var r io.ReadCloser
	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open upload file %s, %v", filename, err)
	}

	// Get the size of the file so that the constraint of Content-Length
	// can be included with the presigned URL. This can be used by the
	// server or client to ensure the content uploaded is of a certain size.
	//
	// These constraints can further be expanded to include things like
	// Content-Type. Additionally constraints such as X-Amz-Content-Sha256
	// header set restricting the content of the file to only the content
	// the client initially made the request with. This prevents the object
	// from being overwritten or used to upload other unintended content.
	stat, err := f.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat file, %s, %v", filename, err)
	}
	defer f.Close()

	req, err := http.NewRequest("PUT", serverURL, f)
	if err != nil {
		return fmt.Errorf("failed to build presigned request, %v", err)
	}

	req.ContentLength = stat.Size()
	req.Header.Add("Content-Length", strconv.Itoa(int(req.ContentLength)))

	// Upload the file contents to S3.
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to do PUT request, %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to put S3 object, %d:%s", resp.StatusCode, resp.Status)
	}

	log.Printf("S3STATUS: %d, %s\n", resp.StatusCode, resp.Status)

	return nil
}

func resourceUploadUrlRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	// svc := meta.(Client).LexBotV2Client

	// urlResp, err := svc.CreateUploadUrl(&lexmodelsv2.CreateUploadUrlInput{})
	// if err != nil {
	// 	return diag.FromErr(err)
	// }

	// d.Set("import_id", d.Get("import_id"))
	// d.Set("upload_url", d.Get("upload_url"))

	// time.Sleep(10 * time.Second)

	return diags
}

func resourceUploadUrlUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	svc := meta.(Client).LexBotV2Client

	urlResp, err := svc.CreateUploadUrl(&lexmodelsv2.CreateUploadUrlInput{})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(aws.StringValue(urlResp.ImportId))
	d.Set("import_id", aws.StringValue(urlResp.ImportId))
	d.Set("upload_url", aws.StringValue(urlResp.UploadUrl))

	// time.Sleep(10 * time.Second)

	err = uploadFile(aws.StringValue(urlResp.UploadUrl), d.Get("file_path").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = svc.StartImport(&lexmodelsv2.StartImportInput{
		ImportId:      urlResp.ImportId,
		MergeStrategy: aws.String(lexmodelsv2.MergeStrategyOverwrite),
		ResourceSpecification: &lexmodelsv2.ImportResourceSpecification{
			BotImportSpecification: &lexmodelsv2.BotImportSpecification{
				BotName:                 aws.String(d.Get("bot_name").(string)),
				DataPrivacy:             &lexmodelsv2.DataPrivacy{ChildDirected: aws.Bool(false)},
				RoleArn:                 aws.String(d.Get("role_arn").(string)),
				IdleSessionTTLInSeconds: aws.Int64(int64(d.Get("idle_session_ttl").(int))),
			},
		},
	})
	if err != nil {
		return diag.FromErr(err)
	}

	respDescr, err := svc.DescribeImport(&lexmodelsv2.DescribeImportInput{
		ImportId: urlResp.ImportId,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	botId := respDescr.ImportedResourceId

	respAliases, err := svc.ListBotAliases(&lexmodelsv2.ListBotAliasesInput{
		BotId: botId,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	if len(respAliases.BotAliasSummaries) == 0 {
		return diag.FromErr(errors.New("no aliases for the bot"))
	}
	aliasId := respAliases.BotAliasSummaries[0].BotAliasId
	region := svc.Config.Region

	stsSvc := meta.(Client).STSClient
	reqAcc, respAcc := stsSvc.GetCallerIdentityRequest(&sts.GetCallerIdentityInput{})
	err = reqAcc.Send()
	if err != nil {
		return diag.FromErr(err)
	}

	aliasArn := fmt.Sprintf("arn:aws:lex:%s:%s:bot-alias/%s/%s",
		aws.StringValue(region),
		aws.StringValue(respAcc.Account),
		aws.StringValue(botId),
		aws.StringValue(aliasId))

	d.Set("bot_alias_arn", aliasArn)

	lambdaArn := d.Get("lambda_arn").(string)
	lambdaVersion := "1.0"

	if len(lambdaArn) > 0 {
		time.Sleep(30 * time.Second)
		reqAlias, respAlias := svc.DescribeBotAliasRequest(&lexmodelsv2.DescribeBotAliasInput{
			BotId:      botId,
			BotAliasId: aliasId,
		})
		err = reqAlias.Send()
		if err != nil {
			log.Println("[DEBUG] calling DescribeBotAliasRequest error")
			return diag.FromErr(err)
		}
		log.Println("[DEBUG] calling DescribeBotAliasRequest success")

		respAlias.BotAliasLocaleSettings = make(map[string]*lexmodelsv2.BotAliasLocaleSettings)
		respAlias.BotAliasLocaleSettings["en_US"] = &lexmodelsv2.BotAliasLocaleSettings{
			Enabled: aws.Bool(true),
			CodeHookSpecification: &lexmodelsv2.CodeHookSpecification{
				LambdaCodeHook: &lexmodelsv2.LambdaCodeHook{
					LambdaARN:                aws.String(lambdaArn),
					CodeHookInterfaceVersion: aws.String(lambdaVersion),
				},
			},
		}
		reqUpdateAlias, respUpdateAlias := svc.UpdateBotAliasRequest((&lexmodelsv2.UpdateBotAliasInput{
			BotAliasId:                respAlias.BotAliasId,
			BotAliasLocaleSettings:    respAlias.BotAliasLocaleSettings,
			BotAliasName:              respAlias.BotAliasName,
			BotId:                     respAlias.BotId,
			BotVersion:                respAlias.BotVersion,
			ConversationLogSettings:   respAlias.ConversationLogSettings,
			Description:               respAlias.Description,
			SentimentAnalysisSettings: respAlias.SentimentAnalysisSettings,
		}))
		err = reqUpdateAlias.Send()
		if err != nil {
			log.Println("[DEBUG] calling UpdateBotAliasRequest error")
			return diag.FromErr(err)
		}
		log.Println("[DEBUG] calling UpdateBotAliasRequest success", respUpdateAlias)
	}

	return diags
}

func resourceUploadUrlDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	// svc := meta.(Client).LexBotV2Client

	// params := &lexmodelsv2.DeleteUploadUrlInput{
	// 	UploadUrlId:            aws.String(d.Id()),
	// 	SkipResourceInUseCheck: aws.Bool(true),
	// }

	// _, err := svc.DeleteUploadUrl(params)
	// if err != nil {
	// 	return diag.FromErr(err)
	// }

	// // d.SetId("") is automatically called assuming delete returns no errors, but
	// // it is added here for explicitness.
	d.SetId("")

	// time.Sleep(10 * time.Second)

	return diags
}
