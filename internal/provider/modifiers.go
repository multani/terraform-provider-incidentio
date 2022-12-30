package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// boolDefaultValueModifier sets a default value for a boolean if none is set
type boolDefaultValueModifier struct {
	DefaultValue bool
}

func boolDefaultValue(value bool) boolDefaultValueModifier {
	return boolDefaultValueModifier{
		DefaultValue: value,
	}
}

func (m boolDefaultValueModifier) Description(ctx context.Context) string {
	return fmt.Sprintf("Default to '%v' if unset", m.DefaultValue)
}

func (m boolDefaultValueModifier) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("Default to `%v` if unset", m.DefaultValue)
}

func (m boolDefaultValueModifier) PlanModifyBool(ctx context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
	if req.ConfigValue.IsNull() {
		resp.PlanValue = types.BoolValue(m.DefaultValue)
	}
}

// int64DefaultValueModifier sets a default value for a string if none is set
type int64DefaultValueModifier struct {
	DefaultValue int64
}

func int64DefaultValue(value int64) int64DefaultValueModifier {
	return int64DefaultValueModifier{
		DefaultValue: value,
	}
}

func (m int64DefaultValueModifier) Description(ctx context.Context) string {
	return fmt.Sprintf("Default to '%v' if unset", m.DefaultValue)
}

func (m int64DefaultValueModifier) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("Default to `%v` if unset", m.DefaultValue)
}

func (m int64DefaultValueModifier) PlanModifyInt64(ctx context.Context, req planmodifier.Int64Request, resp *planmodifier.Int64Response) {
	if req.ConfigValue.IsNull() {
		resp.PlanValue = types.Int64Value(m.DefaultValue)
	}
}

// stringDefaultValueModifier sets a default value for a string if none is set
type stringDefaultValueModifier struct {
	DefaultValue string
}

func stringDefaultValue(value string) stringDefaultValueModifier {
	return stringDefaultValueModifier{
		DefaultValue: value,
	}
}

func (m stringDefaultValueModifier) Description(ctx context.Context) string {
	return fmt.Sprintf("Default to '%v' if unset", m.DefaultValue)
}

func (m stringDefaultValueModifier) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("Default to `%v` if unset", m.DefaultValue)
}

func (m stringDefaultValueModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.ConfigValue.IsNull() {
		resp.PlanValue = types.StringValue(m.DefaultValue)
	}
}
