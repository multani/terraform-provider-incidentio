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
var _ tfsdk.ResourceType = incidentRoleType{}
var _ tfsdk.Resource = incidentRole{}
var _ tfsdk.ResourceWithImportState = incidentRole{}

type incidentRoleType struct{}

func (t incidentRoleType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Configure an incident role",

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Computed:            true,
				MarkdownDescription: "Example identifier",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					tfsdk.UseStateForUnknown(),
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

func (t incidentRoleType) NewResource(ctx context.Context, in tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return incidentRole{
		provider: provider,
	}, diags
}

type incidentRoleData struct {
	Id           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	Required     types.Bool   `tfsdk:"required"`
	Instructions types.String `tfsdk:"instructions"`
	ShortForm    types.String `tfsdk:"short_form"`
}

type incidentRole struct {
	provider provider
}

func (r incidentRole) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var data incidentRoleData

	diags := req.Config.Get(ctx, &data)
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
	response, err := r.provider.client.IncidentRoles().Create(newRole)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create incident role, got error: %s", err))
		return
	}

	data.Id = types.String{Value: response.IncidentRole.Id}
	tflog.Trace(ctx, fmt.Sprintf("created a resource with ID=%s", response.IncidentRole.Id))

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r incidentRole) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var data incidentRoleData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	roleId := data.Id.Value

	response, err := r.provider.client.IncidentRoles().Get(roleId)
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

func (r incidentRole) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
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

	_, err := r.provider.client.IncidentRoles().Update(roleId, updatedRole)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update incident role, got error: %s", err))
		return
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func (r incidentRole) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var data incidentRoleData

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	roleId := data.Id.Value

	err := r.provider.client.IncidentRoles().Delete(roleId)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete incident role, got error: %s", err))
		return
	}
}

func (r incidentRole) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)
}
