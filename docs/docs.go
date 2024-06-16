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
        "/api/v1/executeCode": {
            "post": {
                "description": "根据用户提交的代码和语言执行代码并返回结果",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Code Execution"
                ],
                "summary": "执行代码",
                "parameters": [
                    {
                        "description": "执行代码请求",
                        "name": "executeCodeRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.ExecuteCodeRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "成功响应",
                        "schema": {
                            "$ref": "#/definitions/responses.Response"
                        }
                    },
                    "400": {
                        "description": "错误响应",
                        "schema": {
                            "$ref": "#/definitions/responses.Response"
                        }
                    },
                    "500": {
                        "description": "系统内部错误",
                        "schema": {
                            "$ref": "#/definitions/responses.Response"
                        }
                    }
                }
            }
        },
        "/api/v1/languages": {
            "get": {
                "description": "获取支持的语言列表，只有在该列表中的语言代码才能运行",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Languages"
                ],
                "summary": "获取支持的语言列表",
                "responses": {
                    "200": {
                        "description": "成功响应",
                        "schema": {
                            "$ref": "#/definitions/responses.Response"
                        }
                    },
                    "400": {
                        "description": "错误响应",
                        "schema": {
                            "$ref": "#/definitions/responses.Response"
                        }
                    },
                    "500": {
                        "description": "系统内部错误",
                        "schema": {
                            "$ref": "#/definitions/responses.Response"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "dto.ExecuteCodeRequest": {
            "description": "用户需要提交的代码执行请求参数",
            "type": "object",
            "properties": {
                "code": {
                    "type": "string"
                },
                "inputList": {
                    "description": "运行用例组，每一个元素代表一个用例，例如 1 2\\n",
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "language": {
                    "type": "string"
                }
            }
        },
        "responses.Response": {
            "type": "object",
            "properties": {
                "code": {
                    "description": "错误码",
                    "type": "integer"
                },
                "data": {
                    "description": "返回数据"
                },
                "msg": {
                    "description": "错误描述",
                    "type": "string"
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
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
