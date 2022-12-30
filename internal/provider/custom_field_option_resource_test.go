package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func testAccCustomFieldOptionResourceConfig(field_type string, value string, sort int) string {
	return fmt.Sprintf(`
	resource "incidentio_custom_field" "test" {
		name         = "field1"
		description  = "A description"
		field_type   = "%s"

		required = "always"

		show_before_closure  = true
		show_before_creation = true
		show_before_update   = false
	}

	resource "incidentio_custom_field_option" "test" {
		custom_field_id = incidentio_custom_field.test.id

		value    = "%s"
		sort_key = %d
	}

`, field_type, value, sort)
}

func TestAccCustomFieldOptionResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccCustomFieldOptionResourceConfig("single_select", "test1", 40),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("incidentio_custom_field_option.test", "value", "test1"),
					resource.TestCheckResourceAttr("incidentio_custom_field_option.test", "sort_key", "40"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "incidentio_custom_field_option.test",
				ImportState:       true,
				ImportStateVerify: true,
				// This is not normally necessary, but is here because this
				// example code does not have an actual upstream service.
				// Once the Read method is able to refresh information from
				// the upstream service, this can be removed.
				ImportStateVerifyIgnore: []string{"name"},
			},
			// Update and Read testing
			{
				Config: testAccCustomFieldOptionResourceConfig("single_select", "test2", 42),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("incidentio_custom_field_option.test", "value", "test2"),
					resource.TestCheckResourceAttr("incidentio_custom_field_option.test", "sort_key", "42"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
