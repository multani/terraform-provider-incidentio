package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
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

func (r *SeverityResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Configure a severity",

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Computed:            true,
				MarkdownDescription: "Unique identifier for the severity",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
				Type: types.StringType,
			},
			"name": {
				MarkdownDescription: "Human readable name of the severity",
				Required:            true,
				Type:                types.StringType,
				Validators: []tfsdk.AttributeValidator{
					stringLengthBetween(0, 50),
				},
			},
			"description": {
				MarkdownDescription: "Description of the severity",
				Required:            true,
				Type:                types.StringType,
				Validators: []tfsdk.AttributeValidator{
					stringLengthBetween(0, 1000),
				},
			},
			"rank": {
				MarkdownDescription: "Rank to help sort severities (lower numbers are less severe)",
				Required:            true,
				Type:                types.Int64Type,
			},
		},
	}, nil
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
		Name:        data.Name.Value,
		Description: data.Description.Value,
		Rank:        data.Rank.Value,
	}
	response, err := r.client.Severities().Create(newSeverity)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create severity, got error: %s", err))
		return
	}

	data.Id = types.String{Value: response.Severity.Id}
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

	severityId := data.Id.Value

	response, err := r.client.Severities().Get(severityId)
	if incidentio.IsErrorStatus(err, 404) {
		resp.State.RemoveResource(ctx)
		return
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get severity, got error: %s", err))
		return
	}

	data.Id = types.String{Value: response.Severity.Id}
	data.Name = types.String{Value: response.Severity.Name}
	data.Description = types.String{Value: response.Severity.Description}
	data.Rank = types.Int64{Value: response.Severity.Rank}

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

	severityId := data.Id.Value
	updatedSeverity := incidentio.Severity{
		Name:        data.Name.Value,
		Description: data.Description.Value,
		Rank:        data.Rank.Value,
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

	err := r.client.Severities().Delete(data.Id.Value)
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
