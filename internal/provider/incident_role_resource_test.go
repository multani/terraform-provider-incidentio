package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIncidentRoleResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				VersionConstraint: "3.1.3",
				Source:            "hashicorp/random",
			},
		},
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccIncidentRoleResourceConfig("role 1", false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("incidentio_incident_role.test", "name", "role 1"),
					resource.TestCheckResourceAttr("incidentio_incident_role.test", "required", "false"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "incidentio_incident_role.test",
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
				Config: testAccIncidentRoleResourceConfig("role two", true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("incidentio_incident_role.test", "name", "role two"),
					resource.TestCheckResourceAttr("incidentio_incident_role.test", "required", "true"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccIncidentRoleResourceConfig(name string, required bool) string {
	return fmt.Sprintf(`
	resource "random_string" "short_form" {
		length  = 12
		lower   = true
		upper   = false
		special = false
		number  = false
	}

	resource "incidentio_incident_role" "test" {
		name         = "%s"
		short_form   = random_string.short_form.result
		required     = %v
		description  = "A description"
		instructions = "Some instructions"
	}
`, name, required)
}
