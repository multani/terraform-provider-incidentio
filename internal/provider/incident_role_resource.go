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

func (r *IncidentRoleResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Configure an incident role",

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Computed:            true,
				MarkdownDescription: "Unique identifier for the role",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
				Type: types.StringType,
			},
			"name": {
				MarkdownDescription: "Human readable name of the incident role",
				Required:            true,
				Type:                types.StringType,
			},
			"description": {
				MarkdownDescription: "Describes the purpose of the role",
				Required:            true,
				Type:                types.StringType,
			},
			"required": {
				MarkdownDescription: "Whether incident require this role to be set",
				Required:            true,
				Type:                types.BoolType,
			},
			"instructions": {
				MarkdownDescription: "Provided to whoever is nominated for the role",
				Required:            true,
				Type:                types.StringType,
			},
			"short_form": {
				MarkdownDescription: "Short human readable name for Slack",
				Required:            true,
				Type:                types.StringType,
			},
		},
	}, nil
}

func (r *IncidentRoleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Provider not yet configured
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*incidentio.Client)
}

func (r *IncidentRoleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data incidentRoleData

	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	newRole := incidentio.IncidentRole{
		Name:         data.Name.Value,
		Description:  data.Description.Value,
		Required:     data.Required.Value,
		Instructions: data.Instructions.Value,
		ShortForm:    data.ShortForm.Value,
	}
	response, err := r.client.IncidentRoles().Create(newRole)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create incident role, got error: %s", err))
		return
	}

	data.Id = types.String{Value: response.IncidentRole.Id}
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

	roleId := data.Id.Value

	response, err := r.client.IncidentRoles().Get(roleId)
	if incidentio.IsErrorStatus(err, 404) {
		resp.State.RemoveResource(ctx)
		return
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get incident role, got error: %s", err))
		return
	}

	data.Id = types.String{Value: response.IncidentRole.Id}
	data.Name = types.String{Value: response.IncidentRole.Name}
	data.Description = types.String{Value: response.IncidentRole.Description}
	data.Required = types.Bool{Value: response.IncidentRole.Required}
	data.Instructions = types.String{Value: response.IncidentRole.Instructions}
	data.ShortForm = types.String{Value: response.IncidentRole.ShortForm}

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

	roleId := data.Id.Value
	updatedRole := incidentio.IncidentRole{
		Name:         data.Name.Value,
		Description:  data.Description.Value,
		Required:     data.Required.Value,
		Instructions: data.Instructions.Value,
		ShortForm:    data.ShortForm.Value,
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

	err := r.client.IncidentRoles().Delete(data.Id.Value)
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
