package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/multani/terraform-provider-incidentio/incidentio"
)

type stringLengthBetweenValidator struct {
	Max int
	Min int
}

func stringLengthBetween(minLength int, maxLength int) stringLengthBetweenValidator {
	return stringLengthBetweenValidator{
		Max: maxLength,
		Min: minLength,
	}
}

// Description returns a plain text description of the validator's behavior, suitable for a practitioner to understand its impact.
func (v stringLengthBetweenValidator) Description(ctx context.Context) string {
	return fmt.Sprintf("string length must be between %d and %d", v.Min, v.Max)
}

// MarkdownDescription returns a markdown formatted description of the validator's behavior, suitable for a practitioner to understand its impact.
func (v stringLengthBetweenValidator) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("string length must be between `%d` and `%d`", v.Min, v.Max)
}

// Validate runs the main validation logic of the validator, reading configuration data out of `req` and updating `resp` with diagnostics.
func (v stringLengthBetweenValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	// types.String must be the attr.Value produced by the attr.Type in the schema for this attribute
	// for generic validators, use
	// https://pkg.go.dev/github.com/hashicorp/terraform-plugin-framework/tfsdk#ConvertValue
	// to convert into a known type.
	var str types.String
	diags := tfsdk.ValueAs(ctx, req.ConfigValue, &str)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	if str.IsUnknown() || str.IsNull() {
		return
	}

	strLen := len(str.ValueString())

	if strLen < v.Min || strLen > v.Max {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid String Length",
			fmt.Sprintf("String length must be between %d and %d, got: %d.", v.Min, v.Max, strLen),
		)
		return
	}
}

type fieldTypeValidator struct{}

func isValidCustomFieldFieldType() fieldTypeValidator {
	return fieldTypeValidator{}
}

func (v fieldTypeValidator) Description(ctx context.Context) string {
	return "field type must be one of 'single_select', 'multi_select', 'text', 'link' or 'numeric'"
}

func (v fieldTypeValidator) MarkdownDescription(ctx context.Context) string {
	return "field type must be one of `single_select`, `multi_select`, `text`, `link` or `numeric`"
}

func (v fieldTypeValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	var str types.String
	diags := tfsdk.ValueAs(ctx, req.ConfigValue, &str)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	if str.IsUnknown() || str.IsNull() {
		return
	}

	_, err := incidentio.ParseFieldType(str.ValueString())
	if err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid field type",
			"A field type must be one of 'single_select', 'multi_select', 'text', 'link' or 'numeric'.",
		)
		return
	}
}

type fieldRequirementValidator struct{}

func isValidCustomFieldRequired() fieldRequirementValidator {
	return fieldRequirementValidator{}
}

func (v fieldRequirementValidator) Description(ctx context.Context) string {
	return "Field requirement must be one of 'never', 'before_closure' or 'always'"
}

func (v fieldRequirementValidator) MarkdownDescription(ctx context.Context) string {
	return "Field requirement must be one of `never`, `before_closure` or `always`"
}

func (v fieldRequirementValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	var str types.String
	diags := tfsdk.ValueAs(ctx, req.ConfigValue, &str)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	if str.IsUnknown() || str.IsNull() {
		return
	}

	_, err := incidentio.ParseFieldRequirement(str.ValueString())
	if err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid field requirement",
			"A field requirement must be one of 'never', 'before_closure' or 'always'.",
		)
		return
	}
}
