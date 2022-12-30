package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/multani/terraform-provider-incidentio/incidentio"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &IncidentRoleResource{}
var _ resource.ResourceWithImportState = &IncidentRoleResource{}

type incidentRoleData struct {
	Id           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	Required     types.Bool   `tfsdk:"required"`
	Instructions types.String `tfsdk:"instructions"`
	ShortForm    types.String `tfsdk:"short_form"`
}

type IncidentRoleResource struct {
	// client is the SDK used to communicate with the incident.io service.
	// Resource and DataSource implementations can then make calls using this
	// client.
	client *incidentio.Client
}

func NewIncidentRoleResource() resource.Resource {
	return &IncidentRoleResource{}
}

func (r *IncidentRoleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_incident_role"
}

func (r *IncidentRoleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Configure an incident role",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique identifier for the role",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Human readable name of the incident role",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Describes the purpose of the role",
				Required:            true,
			},
			"required": schema.BoolAttribute{
				MarkdownDescription: "Whether incident require this role to be set",
				Required:            true,
			},
			"instructions": schema.StringAttribute{
				MarkdownDescription: "Provided to whoever is nominated for the role",
				Required:            true,
			},
			"short_form": schema.StringAttribute{
				MarkdownDescription: "Short human readable name for Slack",
				Required:            true,
			},
		},
	}
}

func (r *IncidentRoleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *IncidentRoleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data incidentRoleData

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	newRole := incidentio.IncidentRole{
		Name:         data.Name.ValueString(),
		Description:  data.Description.ValueString(),
		Required:     data.Required.ValueBool(),
		Instructions: data.Instructions.ValueString(),
		ShortForm:    data.ShortForm.ValueString(),
	}
	response, err := r.client.IncidentRoles().Create(newRole)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create incident role, got error: %s", err))
		return
	}

	data.Id = types.StringValue(response.IncidentRole.Id)
	tflog.Trace(ctx, fmt.Sprintf("created a resource with ID=%s", response.IncidentRole.Id))

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r *IncidentRoleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data incidentRoleData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	roleId := data.Id.ValueString()

	response, err := r.client.IncidentRoles().Get(roleId)
	if incidentio.IsErrorStatus(err, 404) {
		resp.State.RemoveResource(ctx)
		return
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get incident role, got error: %s", err))
		return
	}

	data.Id = types.StringValue(response.IncidentRole.Id)
	data.Name = types.StringValue(response.IncidentRole.Name)
	data.Description = types.StringValue(response.IncidentRole.Description)
	data.Required = types.BoolValue(response.IncidentRole.Required)
	data.Instructions = types.StringValue(response.IncidentRole.Instructions)
	data.ShortForm = types.StringValue(response.IncidentRole.ShortForm)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r *IncidentRoleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data incidentRoleData

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	roleId := data.Id.ValueString()
	updatedRole := incidentio.IncidentRole{
		Name:         data.Name.ValueString(),
		Description:  data.Description.ValueString(),
		Required:     data.Required.ValueBool(),
		Instructions: data.Instructions.ValueString(),
		ShortForm:    data.ShortForm.ValueString(),
	}

	_, err := r.client.IncidentRoles().Update(roleId, updatedRole)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update incident role, got error: %s", err))
		return
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r *IncidentRoleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data incidentRoleData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.IncidentRoles().Delete(data.Id.ValueString())
	if incidentio.IsErrorStatus(err, 404) {
		// The resource is already gone.
		return
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete incident role, got error: %s", err))
		return
	}
}

func (r *IncidentRoleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
