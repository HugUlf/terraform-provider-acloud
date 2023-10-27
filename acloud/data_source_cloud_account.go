package acloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/avisi-cloud/go-client/pkg/acloudapi"
)

func dataSourceCloudAccount() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataCloudAccountRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"identity": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"display_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Cloud Account",
			},
			"organisation": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Slug of the Organisation",
			},
			"cloud_provider": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Slug of the Cloud Provider of the Cloud Account",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Returns if the Cloud Account is enabled",
			},
		},
	}
}

func dataCloudAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(acloudapi.Client)
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	org := d.Get("organisation").(string)
	displayName := d.Get("display_name").(string)
	cloudProvider := d.Get("cloud_provider").(string)

	cloudAccount, err := client.FindCloudAccountByName(ctx, org, displayName, cloudProvider)
	if err != nil {
		return diag.FromErr(err)
	}
	if cloudAccount != nil {
		d.SetId(cloudAccount.Identity)
		d.Set("identity", cloudAccount.Identity)
		d.Set("enabled", cloudAccount.Enabled)
		return diags
	}
	return diag.FromErr(fmt.Errorf("cloud account was not found"))
}
