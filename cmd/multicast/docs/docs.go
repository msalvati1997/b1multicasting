// Package docs GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import (
	"bytes"
	"encoding/json"
	"strings"
	"text/template"

	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "email": "soberkoder@gmail.com"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/groups": {
            "get": {
                "description": "Get details of all groups",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "groups"
                ],
                "summary": "Get details of all groups",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.MulticastGroup"
                        }
                    }
                }
            },
            "post": {
                "description": "Create a group from an id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "groups"
                ],
                "summary": "Create a group from an id",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.MulticastGroup"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "main.Member": {
            "type": "object",
            "properties": {
                "address": {
                    "type": "string"
                },
                "member_id": {
                    "type": "string"
                },
                "ready": {
                    "type": "boolean"
                }
            }
        },
        "main.Message": {
            "type": "object",
            "properties": {
                "MessageHeader": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    }
                },
                "Payload": {
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                }
            }
        },
        "main.MulticastGroup": {
            "type": "object",
            "properties": {
                "client_id": {
                    "type": "string"
                },
                "group": {
                    "$ref": "#/definitions/main.MulticastInfo"
                },
                "messages": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/main.Message"
                    }
                }
            }
        },
        "main.MulticastInfo": {
            "type": "object",
            "properties": {
                "members": {
                    "type": "object",
                    "additionalProperties": {
                        "$ref": "#/definitions/main.Member"
                    }
                },
                "multicast_id": {
                    "type": "string"
                },
                "multicast_type": {
                    "type": "string"
                },
                "received_messages": {
                    "type": "integer"
                },
                "status": {
                    "type": "string"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "1.0",
	Host:        "localhost",
	BasePath:    "/",
	Schemes:     []string{},
	Title:       "Orders API",
	Description: "This is a sample service for managing groups multicast",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
		"escape": func(v interface{}) string {
			// escape tabs
			str := strings.Replace(v.(string), "\t", "\\t", -1)
			// replace " with \", and if that results in \\", replace that with \\\"
			str = strings.Replace(str, "\"", "\\\"", -1)
			return strings.Replace(str, "\\\\\"", "\\\\\\\"", -1)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}