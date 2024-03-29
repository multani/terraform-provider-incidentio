package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/multani/terraform-provider-incidentio/incidentio"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &SeverityResource{}
var _ resource.ResourceWithImportState = &SeverityResource{}

type severityData struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Rank        types.Int64  `tfsdk:"rank"`
}

type SeverityResource struct {
	// client is the SDK used to communicate with the incident.io service.
	// Resource and DataSource implementations can then make calls using this
	// client.
	client *incidentio.Client
}

func NewSeverityResource() resource.Resource {
	return &SeverityResource{}
}

func (r *SeverityResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_severity"
}

func (r *SeverityResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Configure a severity",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique identifier for the severity",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Human readable name of the severity",
				Required:            true,
				Validators: []validator.String{
					stringLengthBetween(0, 50),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the severity",
				Required:            true,
				Validators: []validator.String{
					stringLengthBetween(0, 1000),
				},
			},
			"rank": schema.Int64Attribute{
				MarkdownDescription: "Rank to help sort severities (lower numbers are less severe)",
				Required:            true,
			},
		},
	}
}

func (r *SeverityResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Provider not yet configured
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*incidentio.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *incidentio.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *SeverityResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data severityData

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	newSeverity := incidentio.Severity{
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueString(),
		Rank:        data.Rank.ValueInt64(),
	}
	response, err := r.client.Severities().Create(newSeverity)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create severity, got error: %s", err))
		return
	}

	data.Id = types.StringValue(response.Severity.Id)
	tflog.Trace(ctx, fmt.Sprintf("created a resource with ID=%s", response.Severity.Id))

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r *SeverityResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data severityData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	severityId := data.Id.ValueString()

	response, err := r.client.Severities().Get(severityId)
	if incidentio.IsErrorStatus(err, 404) {
		resp.State.RemoveResource(ctx)
		return
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get severity, got error: %s", err))
		return
	}

	data.Id = types.StringValue(response.Severity.Id)
	data.Name = types.StringValue(response.Severity.Name)
	data.Description = types.StringValue(response.Severity.Description)
	data.Rank = types.Int64Value(response.Severity.Rank)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r *SeverityResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data severityData

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	severityId := data.Id.ValueString()
	updatedSeverity := incidentio.Severity{
		Name:        data.Name.ValueString(),
		Description: data.Description.ValueString(),
		Rank:        data.Rank.ValueInt64(),
	}

	_, err := r.client.Severities().Update(severityId, updatedSeverity)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update severity, got error: %s", err))
		return
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r *SeverityResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data severityData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Severities().Delete(data.Id.ValueString())
	if incidentio.IsErrorStatus(err, 404) {
		// The resource is already gone.
		return
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete severity, got error: %s", err))
		return
	}
}

func (r *SeverityResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
