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
        "/send-sms": {
            "post": {
                "description": "Sends an SMS message to a specified phone number",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Send SMS",
                "parameters": [
                    {
                        "description": "Send SMS request",
                        "name": "sms",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.SMSRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "SMS sent successfully",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Invalid request body",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "main.SMSRequest": {
            "description": "SMS request structure with number, message, and sender",
            "type": "object",
            "properties": {
                "date": {
                    "description": "@param date query string false \"Date and time the message will be sent in yyyy-MM-dd HH:mm format. If not provided, the message will be sent as soon as possible\"",
                    "type": "string"
                },
                "message": {
                    "description": "@param message query string true \"The message to be sent\"",
                    "type": "string"
                },
                "number": {
                    "description": "@param number query Use the formatted \"to\" field (either single or array) \"Phone number. The number(s) that will receive the message\""
                },
                "reference": {
                    "description": "@param reference query string false \"Custom reference. A string of max. 255 characters\"",
                    "type": "string"
                },
                "sender": {
                    "description": "@param sender query string true \"Sender. The number or name of the sender. A number can't be longer than 14 characters. A name can't be longer than 11 characters and can't contain special characters or spaces\"",
                    "type": "string"
                },
                "subid": {
                    "description": "@param subid query string false \"ID of a subaccount. If specified, the message will be sent from the subaccount\"",
                    "type": "string"
                },
                "test": {
                    "description": "@param test query bool false \"If true, the system will check all parameters but will not send an SMS message (no credits/balance used)\"",
                    "type": "boolean"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "SMS request structure with number, message, and sender",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
