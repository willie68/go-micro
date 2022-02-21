package apiv1

import "fmt"

const ApiVersion = "1"

var baseURL = fmt.Sprintf("/api/v%s", ApiVersion)

//APIKey the apikey of this service
var APIKey string
