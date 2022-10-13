package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/multani/terraform-provider-incidentio/incidentio"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ provider.Provider = &IncidentIOProvider{}

// IncidentIOProvider satisfies the provider.Provider interface and usually is included
// with all Resource and DataSource implementations.
type IncidentIOProvider struct{}

// providerData can be used to store data from the Terraform configuration.
type providerData struct {
	ApiKey types.String `tfsdk:"api_key"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &IncidentIOProvider{}
	}
}

func (p IncidentIOProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "incidentio"
}

func (p *IncidentIOProvider) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"api_key": {
				MarkdownDescription: "API key. You can also set the `INCIDENT_IO_API_KEY` environment variable instead.",
				Optional:            true,
				Type:                types.StringType,
			},
		},
	}, nil
}

func (p *IncidentIOProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data providerData
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If the upstream provider SDK or HTTP client requires configuration, such
	// as authentication or logging, this is a great opportunity to do so.

	var apiKey string
	if data.ApiKey.Null {
		apiKey = os.Getenv("INCIDENT_IO_API_KEY")
	} else {
		apiKey = data.ApiKey.Value
	}

	client := incidentio.NewClient(apiKey)
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *IncidentIOProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewCustomFieldOptionResource,
		NewCustomFieldResource,
		NewIncidentRoleResource,
		NewSeverityResource,
	}
}

func (p *IncidentIOProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}
