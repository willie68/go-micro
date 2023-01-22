package apiv1

import "fmt"

// APIVersion the version of the implemented api
const APIVersion = "1"

// BaseURL for the routes
var BaseURL = fmt.Sprintf("/api/v%s", APIVersion)

//APIKey the apikey of this service
var APIKey string
