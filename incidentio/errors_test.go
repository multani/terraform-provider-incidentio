package incidentio_test

import (
	"testing"

	"github.com/multani/terraform-provider-incidentio/incidentio"
	"github.com/stretchr/testify/assert"
)

func TestErrorsNew(t *testing.T) {
	data := `
{
	"type":"validation_error",
	"status":422,
	"request_id":"3c9db5ec-36f4-4eed-8bd1-0d9229de7c35",
	"errors":[
	{
		"code":"invalid_value",
		"message":"Shortform must be unique, you already have a role with this shortform!",
		"source":{
			"field":"",
			"pointer":"shortform"
		}
	}
	]
}`

	origErr := incidentio.NewErrors([]byte(data))
	assert.Error(t, origErr)

	err := origErr.(*incidentio.IncidentIOErrorResponse)

	assert.Equal(t, 422, err.Status) // TODO: 422 on the shortform is not documented
	assert.Equal(t,
		"validation_error: invalid_value:Shortform must be unique, you already have a role with this shortform!",
		origErr.Error())

	assert.Equal(t, "", err.Errors[0].Source.Field)
	assert.Equal(t, "shortform", err.Errors[0].Source.Pointer) // TODO: .pointer is not documented
}
