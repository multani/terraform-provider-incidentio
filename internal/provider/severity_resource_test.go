package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSeverityResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		// ExternalProviders: map[string]resource.ExternalProvider{
		// 	"random": {
		// 		VersionConstraint: "3.1.3",
		// 		Source:            "hashicorp/random",
		// 	},
		// },
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccSeverityResourceConfig("sev 1", 10),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("incidentio_severity.test", "name", "sev 1"),
					resource.TestCheckResourceAttr("incidentio_severity.test", "rank", "10"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "incidentio_severity.test",
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
				Config: testAccSeverityResourceConfig("sev 2", 20),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("incidentio_severity.test", "name", "sev 2"),
					resource.TestCheckResourceAttr("incidentio_severity.test", "rank", "20"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccSeverityResourceConfig(name string, rank int) string {
	return fmt.Sprintf(`
	resource "incidentio_severity" "test" {
		name         = "%s"
		description  = "A description"
		rank 		 = %d
	}
`, name, rank)
}
