package provider

import (
	"context"
	"fmt"
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
var _ provider.Provider = &incidentIOProvider{}

// incidentIOProvider satisfies the provider.Provider interface and usually is included
// with all Resource and DataSource implementations.
type incidentIOProvider struct {
	// configured is set to true at the end of the Configure method.
	// This can be used in Resource and DataSource implementations to verify
	// that the provider was previously configured.
	configured bool

	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// providerData can be used to store data from the Terraform configuration.
type providerData struct {
	ApiKey types.String `tfsdk:"api_key"`
}

func (p incidentIOProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "incidentio"
}

func (p *incidentIOProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
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

	p.configured = true
}

func (p *incidentIOProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewCustomFieldOptionResource,
		NewCustomFieldResource,
		NewIncidentRoleResource,
		NewSeverityResource,
	}
}

func (p *incidentIOProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *incidentIOProvider) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
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

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &incidentIOProvider{
			version: version,
		}
	}
}

// convertProviderType is a helper function for NewResource and NewDataSource
// implementations to associate the concrete provider type. Alternatively,
// this helper can be skipped and the provider type can be directly type
// asserted (e.g. provider: in.(*provider)), however using this can prevent
// potential panics.
func convertProviderType(in provider.Provider) (incidentIOProvider, diag.Diagnostics) {
	var diags diag.Diagnostics

	p, ok := in.(*incidentIOProvider)

	if !ok {
		diags.AddError(
			"Unexpected Provider Instance Type",
			fmt.Sprintf("While creating the data source or resource, an unexpected provider type (%T) was received. This is always a bug in the provider code and should be reported to the provider developers.", p),
		)
		return incidentIOProvider{}, diags
	}

	if p == nil {
		diags.AddError(
			"Unexpected Provider Instance Type",
			"While creating the data source or resource, an unexpected empty provider instance was received. This is always a bug in the provider code and should be reported to the provider developers.",
		)
		return incidentIOProvider{}, diags
	}

	return *p, diags
}
