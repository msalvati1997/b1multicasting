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
        "contact": {
            "email": "salvatimartina97@gmail.com"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/groups": {
            "get": {
                "description": "Get Multicast Group",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "groups"
                ],
                "summary": "Get Multicast Group",
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/multicastapp.MulticastInfo"
                        }
                    }
                }
            },
            "post": {
                "description": "Create Multicast Group",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "groups"
                ],
                "summary": "Create Multicast Group",
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/multicastapp.MulticastInfo"
                        }
                    }
                }
            }
        },
        "/groups/:mId": {
            "get": {
                "description": "Get Multicast Group by id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "groups"
                ],
                "summary": "Get Multicast Group by id",
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/multicastapp.MulticastInfo"
                        }
                    }
                }
            },
            "put": {
                "description": "Start multicast by id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "groups"
                ],
                "summary": "Start multicast by id",
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/multicastapp.MulticastInfo"
                        }
                    }
                }
            }
        },
        "/messaging/:mId": {
            "get": {
                "description": "Get Message of Group by id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "groups"
                ],
                "summary": "Get Message of Group by id",
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/multicastapp.Message"
                            }
                        }
                    }
                }
            },
            "post": {
                "description": "Multicast a message to a group G",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "messaging"
                ],
                "summary": "Multicast a message to a group G",
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/multicastapp.Message"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "multicastapp.Member": {
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
        "multicastapp.Message": {
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
        "multicastapp.MulticastInfo": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                },
                "members": {
                    "type": "object",
                    "additionalProperties": {
                        "$ref": "#/definitions/multicastapp.Member"
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
	Host:        "localhost:8080",
	BasePath:    "",
	Schemes:     []string{},
	Title:       "MULTICAST API",
	Description: "This is a multicast API",
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
	swag.Register("swagger", &s{})
}
