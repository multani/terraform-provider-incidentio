package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func testAccCustomFieldResourceConfig(name string, required string, field_type string) string {
	return fmt.Sprintf(`
	resource "incidentio_custom_field" "test" {
		name         = "%s"
		description  = "A description"

		required = "%s"

		show_before_closure  = true
		show_before_creation = true

		field_type = "%s"
	}
`, name, required, field_type)
}

func TestAccCustomFieldResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccCustomFieldResourceConfig("field1", "always", "text"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("incidentio_custom_field.test", "name", "field1"),
					resource.TestCheckResourceAttr("incidentio_custom_field.test", "required", "always"),
					resource.TestCheckResourceAttr("incidentio_custom_field.test", "field_type", "text"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "incidentio_custom_field.test",
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
				Config: testAccCustomFieldResourceConfig("field2", "never", "text"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("incidentio_custom_field.test", "name", "field2"),
					resource.TestCheckResourceAttr("incidentio_custom_field.test", "required", "never"),
					resource.TestCheckResourceAttr("incidentio_custom_field.test", "field_type", "text"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
