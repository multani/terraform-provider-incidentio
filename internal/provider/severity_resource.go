package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/multani/terraform-provider-incidentio/incidentio"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ tfsdk.ResourceType = severityType{}
var _ tfsdk.Resource = severity{}
var _ tfsdk.ResourceWithImportState = severity{}

type severityType struct{}

func (t severityType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Configure a severity",

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Computed:            true,
				MarkdownDescription: "Unique identifier for the severity",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
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

func (t severityType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return severity{
		provider: provider,
	}, diags
}

type severityData struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Rank        types.Int64  `tfsdk:"rank"`
}

type severity struct {
	provider provider
}

func (r severity) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
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
	response, err := r.provider.client.Severities().Create(newSeverity)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create severity, got error: %s", err))
		return
	}

	data.Id = types.String{Value: response.Severity.Id}
	tflog.Trace(ctx, fmt.Sprintf("created a resource with ID=%s", response.Severity.Id))

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r severity) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var data severityData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	severityId := data.Id.Value

	response, err := r.provider.client.Severities().Get(severityId)
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

func (r severity) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
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

	_, err := r.provider.client.Severities().Update(severityId, updatedSeverity)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update severity, got error: %s", err))
		return
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r severity) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var data severityData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.provider.client.Severities().Delete(data.Id.Value)
	if incidentio.IsErrorStatus(err, 404) {
		// The resource is already gone.
		return
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete severity, got error: %s", err))
		return
	}
}

func (r severity) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)
}
