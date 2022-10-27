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
var _ resource.Resource = &CustomFieldOptionResource{}
var _ resource.ResourceWithImportState = &CustomFieldOptionResource{}

type customFieldOptionData struct {
	Id            types.String `tfsdk:"id"`
	CustomFieldId types.String `tfsdk:"custom_field_id"`
	Value         types.String `tfsdk:"value"`
	SortKey       types.Int64  `tfsdk:"sort_key"`
}

type CustomFieldOptionResource struct {
	// client is the SDK used to communicate with the incident.io service.
	// Resource and DataSource implementations can then make calls using this
	// client.
	client *incidentio.Client
}

func NewCustomFieldOptionResource() resource.Resource {
	return &CustomFieldOptionResource{}
}

func (r *CustomFieldOptionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_custom_field_option"
}

func (r *CustomFieldOptionResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Configure a custom field option",

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:                types.StringType,
				Computed:            true,
				MarkdownDescription: "Unique identifier for the custom field option",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
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

func (r *CustomFieldOptionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CustomFieldOptionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
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
	response, err := r.client.CustomFieldOptions().Create(newCustomFieldOption)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create custom field option, got error: %s", err))
		return
	}

	data.Id = types.StringValue(response.CustomFieldOption.Id)
	tflog.Trace(ctx, fmt.Sprintf("created a resource with ID=%s", response.CustomFieldOption.Id))

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r *CustomFieldOptionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data customFieldOptionData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	id := data.Id.Value

	response, err := r.client.CustomFieldOptions().Get(id)
	if incidentio.IsErrorStatus(err, 404) {
		resp.State.RemoveResource(ctx)
		return
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get custom field option, got error: %s", err))
		return
	}

	data.Id = types.StringValue(response.CustomFieldOption.Id)
	data.CustomFieldId = types.StringValue(response.CustomFieldOption.CustomFieldId)
	data.Value = types.StringValue(response.CustomFieldOption.Value)
	data.SortKey = types.Int64Value(response.CustomFieldOption.SortKey)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r *CustomFieldOptionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
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

	_, err := r.client.CustomFieldOptions().Update(id, updatedCFO)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update custom field option, got error: %s", err))
		return
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r *CustomFieldOptionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data customFieldOptionData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.CustomFieldOptions().Delete(data.Id.Value)
	if incidentio.IsErrorStatus(err, 404) {
		// The resource is already gone.
		return
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete custom field option, got error: %s", err))
		return
	}
}

func (r *CustomFieldOptionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
