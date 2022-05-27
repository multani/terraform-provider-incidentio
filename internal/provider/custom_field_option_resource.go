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
var _ tfsdk.ResourceType = customFieldOptionType{}
var _ tfsdk.Resource = customFieldOption{}
var _ tfsdk.ResourceWithImportState = customFieldOption{}

type customFieldOptionType struct{}

func (t customFieldOptionType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Configure a custom field option",

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:                types.StringType,
				Computed:            true,
				MarkdownDescription: "Unique identifier for the custom field option",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
				},
			},
			"custom_field_id": {
				Type:                types.StringType,
				MarkdownDescription: "ID of the custom field this option belongs to",
				Required:            true,
			},
			"value": {
				Type:                types.StringType,
				MarkdownDescription: "Human readable name for the custom field option",
				Required:            true,
			},
			// TODO: make this optional with a default value
			"sort_key": {
				Type:                types.Int64Type,
				MarkdownDescription: "Sort key used to order the custom field options correctly",
				Required:            true,
			},
		},
	}, nil
}

func (t customFieldOptionType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return customFieldOption{
		provider: provider,
	}, diags
}

type customFieldOptionData struct {
	Id            types.String `tfsdk:"id"`
	CustomFieldId types.String `tfsdk:"custom_field_id"`
	Value         types.String `tfsdk:"value"`
	SortKey       types.Int64  `tfsdk:"sort_key"`
}

type customFieldOption struct {
	provider provider
}

func (r customFieldOption) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var data customFieldOptionData

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	newCustomFieldOption := incidentio.CustomFieldOption{
		CustomFieldId: data.CustomFieldId.Value,
		Value:         data.Value.Value,
		SortKey:       data.SortKey.Value,
	}
	response, err := r.provider.client.CustomFieldOptions().Create(newCustomFieldOption)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create custom field option, got error: %s", err))
		return
	}

	data.Id = types.String{Value: response.CustomFieldOption.Id}
	tflog.Trace(ctx, fmt.Sprintf("created a resource with ID=%s", response.CustomFieldOption.Id))

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r customFieldOption) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var data customFieldOptionData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	id := data.Id.Value

	response, err := r.provider.client.CustomFieldOptions().Get(id)
	if incidentio.IsErrorStatus(err, 404) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get custom field option, got error: %s", err))
		return
	}

	data.Id = types.String{Value: response.CustomFieldOption.Id}
	data.CustomFieldId = types.String{Value: response.CustomFieldOption.CustomFieldId}
	data.Value = types.String{Value: response.CustomFieldOption.Value}
	data.SortKey = types.Int64{Value: response.CustomFieldOption.SortKey}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r customFieldOption) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var data customFieldOptionData

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	id := data.Id.Value
	updatedCFO := incidentio.CustomFieldOption{
		CustomFieldId: data.CustomFieldId.Value,
		Value:         data.Value.Value,
		SortKey:       data.SortKey.Value,
	}

	_, err := r.provider.client.CustomFieldOptions().Update(id, updatedCFO)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update custom field option, got error: %s", err))
		return
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r customFieldOption) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var data customFieldOptionData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	severityId := data.Id.Value

	err := r.provider.client.CustomFieldOptions().Delete(severityId)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete custom field option, got error: %s", err))
		return
	}
}

func (r customFieldOption) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)
}
