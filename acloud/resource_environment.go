package acloud

import (
	"context"
	"fmt"
	"strconv"

	"github.com/avisi-cloud/go-client/pkg/acloudapi"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceEnvironment() *schema.Resource {
	return &schema.Resource{
		Description:   "Create an environment",
		CreateContext: resourceEnvironmentCreate,
		ReadContext:   resourceEnvironmentRead,
		UpdateContext: resourceEnvironmentUpdate,
		DeleteContext: resourceEnvironmentDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"organisation": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Slug of the Organisation",
			},
			"organisation_slug": {
				Type:       schema.TypeString,
				Deprecated: "replaced by organisation",
				Optional:   true,
				Default:    nil,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Environment",
			},
			"slug": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"purpose": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Purpose of the Environment",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Type of the Environment. Available options: production, staging, development, demo, other",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A human readable description about the environment",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceEnvironmentCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(acloudapi.Client)
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	createEnvironment := acloudapi.CreateEnvironment{
		Name:        d.Get("name").(string),
		Purpose:     d.Get("purpose").(string),
		Type:        d.Get("type").(string),
		Description: d.Get("description").(string),
	}

	org := getStringAttributeWithLegacyName(d, "organisation", "organisation_slug")

	environment, err := client.CreateEnvironment(ctx, createEnvironment, org)

	if err != nil {
		return diag.FromErr(err)
	}
	if environment != nil {
		d.SetId(strconv.Itoa(environment.ID))
		d.Set("slug", environment.Slug)
		return diags
	}
	return resourceEnvironmentRead(ctx, d, m)
}

func resourceEnvironmentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(acloudapi.Client)
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	org := getStringAttributeWithLegacyName(d, "organisation", "organisation_slug")
	slug := d.Get("slug").(string)
	environment, err := client.GetEnvironment(ctx, org, slug)
	if err != nil {
		return diag.FromErr(err)
	}
	if environment == nil {
		return diag.FromErr(fmt.Errorf("environment was not found"))
	}

	d.SetId(strconv.Itoa(environment.ID))
	d.Set("name", environment.Name)
	d.Set("slug", environment.Slug)
	d.Set("purpose", environment.Purpose)
	d.Set("type", environment.Type)
	d.Set("description", environment.Description)

	return diags

}

func resourceEnvironmentUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(acloudapi.Client)
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	updateEnvironment := acloudapi.UpdateEnvironment{
		Name:        d.Get("name").(string),
		Purpose:     d.Get("purpose").(string),
		Type:        d.Get("type").(string),
		Description: d.Get("description").(string),
	}

	org := getStringAttributeWithLegacyName(d, "organisation", "organisation_slug")
	env := d.Get("slug").(string)

	environment, err := client.UpdateEnvironment(ctx, updateEnvironment, org, env)
	if err != nil {
		return diag.FromErr(err)
	}
	if environment != nil {
		d.Set("name", environment.Name)
		d.Set("purpose", environment.Purpose)
		d.Set("type", environment.Type)
		d.Set("description", environment.Description)
		d.Set("slug", environment.Slug)
		return diags
	}

	return resourceEnvironmentRead(ctx, d, m)
}

func resourceEnvironmentDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(acloudapi.Client)
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	org := getStringAttributeWithLegacyName(d, "organisation", "organisation_slug")
	slug := d.Get("slug").(string)

	err := client.DeleteEnvironment(ctx, org, slug)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
