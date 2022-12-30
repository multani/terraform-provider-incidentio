resource "incidentio_custom_field" "test" {
  name        = "Affected Team"
  description = "The team which was responsible for resolving this incident."

  #required = "never" # never, before_closure

  show_before_closure  = true
  show_before_creation = true

  field_type = "multi_select"

  #condition {
  #operation = "one_of"
  # subject = "incident.severity"

  #option {
  #value = "major"
  #}
  #}
}

locals {
  test_options = [
    "test1",
    "test2",
    "test3",
  ]
}

resource "incidentio_custom_field_option" "test" {
  for_each = {
    for key, value in zipmap(
      local.test_options,
      range(length(local.test_options)),
    ) : key => value
  }

  custom_field_id = incidentio_custom_field.test.id

  value = each.key
  #sort_key = each.value
}
