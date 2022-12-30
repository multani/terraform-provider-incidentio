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
var _ resource.Resource = &CustomFieldResource{}
var _ resource.ResourceWithImportState = &CustomFieldResource{}

type customField struct {
	Id                 types.String `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	Description        types.String `tfsdk:"description"`
	Required           types.String `tfsdk:"required"`
	ShowBeforeClosure  types.Bool   `tfsdk:"show_before_closure"`
	ShowBeforeCreation types.Bool   `tfsdk:"show_before_creation"`
	ShowBeforeUpdate   types.Bool   `tfsdk:"show_before_update"`
	FieldType          types.String `tfsdk:"field_type"`
}

type CustomFieldResource struct {
	// client is the SDK used to communicate with the incident.io service.
	// Resource and DataSource implementations can then make calls using this
	// client.
	client *incidentio.Client
}

func NewCustomFieldResource() resource.Resource {
	return &CustomFieldResource{}
}

func (r *CustomFieldResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_custom_field"
}

func (r *CustomFieldResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Configure a custom field",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique identifier for the custom field",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Human readable name of the custom field",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the custom field",
				Required:            true,
			},
			"field_type": schema.StringAttribute{
				MarkdownDescription: "The type of the custom field",
				Required:            true,
				Validators: []validator.String{
					isValidCustomFieldFieldType(),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"required": schema.StringAttribute{
				MarkdownDescription: "When this custom field must be set during the incident lifecycle. " +
					"Must be one of `never`, `before_closure` or `always`.",
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					isValidCustomFieldRequired(),
				},
				PlanModifiers: []planmodifier.String{
					stringDefaultValue("always"),
				},
			},
			"show_before_closure": schema.BoolAttribute{
				MarkdownDescription: "Whether a custom field should be shown in the incident close modal. If this custom field is required before closure, but no value has been set for it, the field will be shown in the closure modal whatever the value of this setting.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolDefaultValue(true),
				},
			},
			"show_before_creation": schema.BoolAttribute{
				MarkdownDescription: "Whether a custom field should be shown in the incident creation modal. This must be true if the field is always required.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolDefaultValue(true),
				},
			},
			"show_before_update": schema.BoolAttribute{
				MarkdownDescription: "Whether a custom field should be shown in the incident update modal.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolDefaultValue(true),
				},
			},
		},
	}
}

func (r *CustomFieldResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CustomFieldResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data customField

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	newCF := incidentio.CustomField{
		Name:               data.Name.ValueString(),
		Description:        data.Description.ValueString(),
		Required:           incidentio.FieldRequirement(data.Required.ValueString()),
		ShowBeforeClosure:  data.ShowBeforeClosure.ValueBool(),
		ShowBeforeCreation: data.ShowBeforeCreation.ValueBool(),
		ShowBeforeUpdate:   data.ShowBeforeUpdate.ValueBool(),
		FieldType:          incidentio.FieldType(data.FieldType.ValueString()),
	}

	response, err := r.client.CustomFields().Create(newCF)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create custom field, got error: %s", err))
		return
	}

	data.Id = types.StringValue(response.CustomField.Id)
	tflog.Trace(ctx, fmt.Sprintf("created a resource with ID=%s", response.CustomField.Id))

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r *CustomFieldResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data customField

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	id := data.Id.ValueString()

	response, err := r.client.CustomFields().Get(id)
	if incidentio.IsErrorStatus(err, 404) {
		resp.State.RemoveResource(ctx)
		return
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get custom field option, got error: %s", err))
		return
	}

	data.Id = types.StringValue(response.CustomField.Id)
	data.Name = types.StringValue(response.CustomField.Name)
	data.Description = types.StringValue(response.CustomField.Description)
	data.Required = types.StringValue(string(response.CustomField.Required))
	data.ShowBeforeClosure = types.BoolValue(response.CustomField.ShowBeforeClosure)
	data.ShowBeforeCreation = types.BoolValue(response.CustomField.ShowBeforeCreation)
	data.ShowBeforeUpdate = types.BoolValue(response.CustomField.ShowBeforeUpdate)
	data.FieldType = types.StringValue(string(response.CustomField.FieldType))

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r *CustomFieldResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data customField

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	cfId := data.Id.ValueString()

	updatedCF := incidentio.CustomField{
		Name:               data.Name.ValueString(),
		Description:        data.Description.ValueString(),
		Required:           incidentio.FieldRequirement(data.Required.ValueString()),
		ShowBeforeClosure:  data.ShowBeforeClosure.ValueBool(),
		ShowBeforeCreation: data.ShowBeforeCreation.ValueBool(),
		ShowBeforeUpdate:   data.ShowBeforeUpdate.ValueBool(),
		FieldType:          incidentio.FieldType(data.FieldType.ValueString()),
	}

	_, err := r.client.CustomFields().Update(cfId, updatedCF)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update custom field, got error: %s", err))
		return
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r *CustomFieldResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data customField

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.CustomFields().Delete(data.Id.ValueString())
	if incidentio.IsErrorStatus(err, 404) {
		// The resource is already gone.
		return
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete custom field, got error: %s", err))
		return
	}
}

func (r *CustomFieldResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
