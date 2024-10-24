// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Anton",
            "email": "ap363402@gmail.com"
        },
        "license": {
            "name": "Apache 2.0"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/metric": {
            "get": {
                "description": "Checks that the server is up and running",
                "tags": [
                    "Heartbeat"
                ],
                "summary": "Heartbeat",
                "responses": {
                    "204": {
                        "description": "No Content"
                    }
                }
            }
        },
        "/operations": {
            "get": {
                "description": "Retrieves a list of operations with support for filtering and sorting.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Operations"
                ],
                "summary": "Get operations",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User UUID",
                        "name": "user_uuid",
                        "in": "path"
                    },
                    {
                        "type": "string",
                        "description": "Category name (supports operators: substr)",
                        "name": "category_name",
                        "in": "path"
                    },
                    {
                        "type": "string",
                        "description": "Category type",
                        "name": "type",
                        "in": "path"
                    },
                    {
                        "type": "string",
                        "description": "Category ID",
                        "name": "category_id",
                        "in": "path"
                    },
                    {
                        "type": "string",
                        "description": "Description (supports operators: substr)",
                        "name": "description",
                        "in": "path"
                    },
                    {
                        "type": "string",
                        "description": "Money sum (supports operators: eq, neq, lt, lte, gt, gte, between)",
                        "name": "money_sum",
                        "in": "path"
                    },
                    {
                        "type": "string",
                        "description": "Date and time of operation (supports operators: eq, between; format: yyyy-mm-dd)",
                        "name": "date_time",
                        "in": "path"
                    },
                    {
                        "type": "string",
                        "description": "Field to sort by (money_sum, date_time, description)",
                        "name": "sort_by",
                        "in": "path"
                    },
                    {
                        "type": "string",
                        "description": "Sort order (asc, desc)",
                        "name": "sort_order",
                        "in": "path"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "List of operations",
                        "schema": {
                            "$ref": "#/definitions/entity.Report"
                        }
                    },
                    "400": {
                        "description": "Validation error in filter or sort parameters",
                        "schema": {
                            "$ref": "#/definitions/apperror.AppError"
                        }
                    },
                    "418": {
                        "description": "Something wrong with application logic",
                        "schema": {
                            "$ref": "#/definitions/apperror.AppError"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/apperror.AppError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "apperror.AppError": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "string"
                },
                "developer_message": {
                    "type": "string"
                },
                "fields": {
                    "$ref": "#/definitions/apperror.ErrorFields"
                },
                "message": {
                    "type": "string"
                },
                "params": {
                    "$ref": "#/definitions/apperror.ErrorParams"
                }
            }
        },
        "apperror.ErrorFields": {
            "type": "object",
            "additionalProperties": {
                "type": "string"
            }
        },
        "apperror.ErrorParams": {
            "type": "object",
            "additionalProperties": {
                "type": "string"
            }
        },
        "entity.Operation": {
            "type": "object",
            "properties": {
                "category_uuid": {
                    "type": "string"
                },
                "date_time": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "money_sum": {
                    "type": "number"
                },
                "uuid": {
                    "type": "string"
                }
            }
        },
        "entity.Report": {
            "type": "object",
            "properties": {
                "operations": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/entity.Operation"
                    }
                },
                "total_money_sum": {
                    "type": "number"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:10003",
	BasePath:         "/api",
	Schemes:          []string{},
	Title:            "Stats-service API",
	Description:      "Statistics service for finance-manager application",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
