// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

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
        "/users": {
            "get": {
                "description": "Get users with optional filtering by name and age",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Get users with filtering",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Filter by name",
                        "name": "name",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Minimum age",
                        "name": "min_age",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Maximum age",
                        "name": "max_age",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Number of records to return",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Number of records to skip",
                        "name": "skip",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/server.UserResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    },
                    "429": {
                        "description": "Too Many Requests",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/server.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "server.ErrorResponse": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "server.UserResponse": {
            "type": "object",
            "properties": {
                "data": {},
                "message": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:3000",
	BasePath:         "/",
	Schemes:          []string{"http"},
	Title:            "Hashtechy API",
	Description:      "This is a sample server for managing users with encryption and caching",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
