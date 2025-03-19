// Package api Code generated by swaggo/swag. DO NOT EDIT
package api

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/addresses": {
            "get": {
                "security": [
                    {
                        "api_key": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "addresses"
                ],
                "summary": "getting all addresses",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Tenant",
                        "name": "tenant",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "response with list of addresses as json",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/pmodel.Address"
                            }
                        }
                    },
                    "400": {
                        "description": "client error information as json",
                        "schema": {
                            "$ref": "#/definitions/github_com_willie68_go-micro_internal_serror.Serr"
                        }
                    },
                    "500": {
                        "description": "server error information as json",
                        "schema": {
                            "$ref": "#/definitions/github_com_willie68_go-micro_internal_serror.Serr"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "api_key": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "addresses"
                ],
                "summary": "Create a new address",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Tenant",
                        "name": "tenant",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "address to be added",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/pmodel.Address"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "tenant",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "client error information as json",
                        "schema": {
                            "$ref": "#/definitions/github_com_willie68_go-micro_internal_serror.Serr"
                        }
                    },
                    "500": {
                        "description": "server error information as json",
                        "schema": {
                            "$ref": "#/definitions/github_com_willie68_go-micro_internal_serror.Serr"
                        }
                    }
                }
            }
        },
        "/addresses/{id}": {
            "get": {
                "security": [
                    {
                        "api_key": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "addresses"
                ],
                "summary": "getting one address",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Tenant",
                        "name": "tenant",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "response with the address with id as json",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/pmodel.Address"
                            }
                        }
                    },
                    "400": {
                        "description": "client error information as json",
                        "schema": {
                            "$ref": "#/definitions/github_com_willie68_go-micro_internal_serror.Serr"
                        }
                    },
                    "500": {
                        "description": "server error information as json",
                        "schema": {
                            "$ref": "#/definitions/github_com_willie68_go-micro_internal_serror.Serr"
                        }
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "api_key": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "addresses"
                ],
                "summary": "Create a new address",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Tenant",
                        "name": "tenant",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "Add store",
                        "name": "payload",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "tenant",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "client error information as json",
                        "schema": {
                            "$ref": "#/definitions/github_com_willie68_go-micro_internal_serror.Serr"
                        }
                    },
                    "500": {
                        "description": "server error information as json",
                        "schema": {
                            "$ref": "#/definitions/github_com_willie68_go-micro_internal_serror.Serr"
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "api_key": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "addresses"
                ],
                "summary": "Delete a address",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Tenant",
                        "name": "tenant",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok"
                    },
                    "400": {
                        "description": "client error information as json",
                        "schema": {
                            "$ref": "#/definitions/github_com_willie68_go-micro_internal_serror.Serr"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "github_com_willie68_go-micro_internal_serror.Serr": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "key": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                },
                "origin": {
                    "type": "string"
                },
                "service": {
                    "type": "string"
                }
            }
        },
        "pmodel.Address": {
            "type": "object",
            "properties": {
                "city": {
                    "type": "string"
                },
                "firstname": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "state": {
                    "type": "string"
                },
                "street": {
                    "type": "string"
                },
                "zip_code": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "GoMicro service API",
	Description:      "The GoMicro service is a template for microservices written in go.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
