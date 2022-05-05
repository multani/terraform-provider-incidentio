resource "incidentio_incident_role" "lead" {
  #count      = 0
  name       = "Incident Lead (test)"
  short_form = "lead6" # lead is reserved
  required   = false

  description = "The person currently coordinating the incident, tasked with driving it to resolution and ensuring clear internal and external communication with stakeholders and customers."
  #description = "test123"

  instructions = <<EOF
- Make sure it’s clear who is doing what, and that people are working together effectively
- Ensure everybody has what they need, and any blockers are flagged quickly
- Provide regular, clear updates for stakeholders to let them know what’s happening
EOF
}
