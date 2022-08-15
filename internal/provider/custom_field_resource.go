package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/multani/terraform-provider-incidentio/incidentio"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ provider.ResourceType = customFieldType{}
var _ resource.Resource = resourceCustomField{}
var _ resource.ResourceWithImportState = resourceCustomField{}

type customFieldType struct{}

func (t customFieldType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Configure a custom field",

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Computed:            true,
				MarkdownDescription: "Unique identifier for the custom field",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
				Type: types.StringType,
			},
			"name": {
				MarkdownDescription: "Human readable name of the custom field",
				Required:            true,
				Type:                types.StringType,
			},
			"description": {
				MarkdownDescription: "Description of the custom field",
				Required:            true,
				Type:                types.StringType,
			},
			"field_type": {
				MarkdownDescription: "The type of the custom field",
				Required:            true,
				Type:                types.StringType,
				Validators: []tfsdk.AttributeValidator{
					isValidCustomFieldFieldType(),
				},
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
				},
			},
			// TODO: make this optional with a default value
			"required": {
				MarkdownDescription: "When this custom field must be set during the incident lifecycle. " +
					"Must be one of `never`, `before_closure` or `always`.",
				Required: true,
				Type:     types.StringType,
				Validators: []tfsdk.AttributeValidator{
					isValidCustomFieldRequired(),
				},
			},
			// TODO: make this optional with a default value
			"show_before_closure": {
				MarkdownDescription: "Whether a custom field should be shown in the incident close modal. If this custom field is required before closure, but no value has been set for it, the field will be shown in the closure modal whatever the value of this setting.",
				Required:            true,
				Type:                types.BoolType,
			},
			// TODO: make this optional with a default value
			"show_before_creation": {
				MarkdownDescription: "Whether a custom field should be shown in the incident creation modal. This must be true if the field is always required.",
				Required:            true,
				Type:                types.BoolType,
			},
		},
	}, nil
}

func (t customFieldType) NewResource(ctx context.Context, in provider.Provider) (resource.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return resourceCustomField{
		provider: provider,
	}, diags
}

type customField struct {
	Id                 types.String `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	Description        types.String `tfsdk:"description"`
	Required           types.String `tfsdk:"required"`
	ShowBeforeClosure  types.Bool   `tfsdk:"show_before_closure"`
	ShowBeforeCreation types.Bool   `tfsdk:"show_before_creation"`
	FieldType          types.String `tfsdk:"field_type"`
}

type resourceCustomField struct {
	provider incidentIOProvider
}

func (r resourceCustomField) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data customField

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	newCF := incidentio.CustomField{
		Name:               data.Name.Value,
		Description:        data.Description.Value,
		Required:           incidentio.FieldRequirement(data.Required.Value),
		ShowBeforeClosure:  data.ShowBeforeClosure.Value,
		ShowBeforeCreation: data.ShowBeforeCreation.Value,
		FieldType:          incidentio.FieldType(data.FieldType.Value),
	}

	response, err := r.provider.client.CustomFields().Create(newCF)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create custom field, got error: %s", err))
		return
	}

	data.Id = types.String{Value: response.CustomField.Id}
	tflog.Trace(ctx, fmt.Sprintf("created a resource with ID=%s", response.CustomField.Id))

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r resourceCustomField) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data customField

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	id := data.Id.Value

	response, err := r.provider.client.CustomFields().Get(id)
	if incidentio.IsErrorStatus(err, 404) {
		resp.State.RemoveResource(ctx)
		return
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get custom field option, got error: %s", err))
		return
	}

	data.Id = types.String{Value: response.CustomField.Id}
	data.Name = types.String{Value: response.CustomField.Name}
	data.Description = types.String{Value: response.CustomField.Description}
	data.Required = types.String{Value: string(response.CustomField.Required)}
	data.ShowBeforeClosure = types.Bool{Value: response.CustomField.ShowBeforeClosure}
	data.ShowBeforeCreation = types.Bool{Value: response.CustomField.ShowBeforeCreation}
	data.FieldType = types.String{Value: string(response.CustomField.FieldType)}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r resourceCustomField) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data customField

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	cfId := data.Id.Value

	updatedCF := incidentio.CustomField{
		Name:               data.Name.Value,
		Description:        data.Description.Value,
		Required:           incidentio.FieldRequirement(data.Required.Value),
		ShowBeforeClosure:  data.ShowBeforeClosure.Value,
		ShowBeforeCreation: data.ShowBeforeCreation.Value,
		FieldType:          incidentio.FieldType(data.FieldType.Value),
	}

	_, err := r.provider.client.CustomFields().Update(cfId, updatedCF)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update custom field, got error: %s", err))
		return
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r resourceCustomField) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data customField

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.provider.client.CustomFields().Delete(data.Id.Value)
	if incidentio.IsErrorStatus(err, 404) {
		// The resource is already gone.
		return
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete custom field, got error: %s", err))
		return
	}
}

func (r resourceCustomField) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
