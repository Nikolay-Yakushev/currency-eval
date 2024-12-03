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
        "/currencies": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Currencies"
                ],
                "parameters": [
                    {
                        "description": "Request body",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.RequestCurrencyPairDTO"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.ResponseCurrencyPairDTO"
                        }
                    }
                }
            }
        },
        "/currencies_with_date": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Currencies"
                ],
                "parameters": [
                    {
                        "description": "Request body",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.RequestCurrencyByDateDTO"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.ResponseCurrencyByDateDTO"
                        }
                    }
                }
            }
        },
        "/ping": {
            "get": {
                "description": "Responds with a \"pong\" message",
                "tags": [
                    "Ping"
                ],
                "summary": "Ping the server",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "dto.RequestCurrencyByDateDTO": {
            "type": "object",
            "properties": {
                "baseCurrency": {
                    "type": "string"
                },
                "effectiveDate": {
                    "type": "string"
                }
            }
        },
        "dto.RequestCurrencyPairDTO": {
            "type": "object",
            "properties": {
                "baseCurrency": {
                    "description": "relative to which currency rates should be calculated",
                    "type": "string"
                },
                "targetCurrency": {
                    "type": "string"
                }
            }
        },
        "dto.ResponseCurrencyByDateDTO": {
            "type": "object",
            "properties": {
                "baseCurrency": {
                    "type": "string"
                },
                "baseCurrencyValue": {
                    "type": "number"
                },
                "currencies": {
                    "description": "{EUR: 1.23} // value relative to BaseCurrency value",
                    "type": "object",
                    "additionalProperties": {
                        "type": "number"
                    }
                },
                "updatedAt": {
                    "type": "string"
                }
            }
        },
        "dto.ResponseCurrencyPairDTO": {
            "type": "object",
            "properties": {
                "baseCurrency": {
                    "type": "string"
                },
                "baseCurrencyValue": {
                    "type": "number"
                },
                "targetCurrency": {
                    "type": "string"
                },
                "targetCurrencyValue": {
                    "type": "number"
                },
                "updateAt": {
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
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Currency API",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}