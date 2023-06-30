// Code generated by swaggo/swag. DO NOT EDIT.

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
        "/images/download/{imageName}": {
            "get": {
                "description": "Downloads the specified image from the server",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/octet-stream"
                ],
                "tags": [
                    "Images"
                ],
                "summary": "Download an image",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Image filename",
                        "name": "imageName",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Image file",
                        "schema": {
                            "type": "file"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {}
                    }
                }
            }
        },
        "/images/upload": {
            "post": {
                "description": "Uploads an image to the server",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Images"
                ],
                "summary": "Upload an image",
                "parameters": [
                    {
                        "type": "file",
                        "description": "Image file",
                        "name": "image",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {}
                    }
                }
            }
        },
        "/login": {
            "post": {
                "description": "Authenticates a user with the provided username and password",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "User login",
                "parameters": [
                    {
                        "description": "userRequest value (model.PersonRequest)",
                        "name": "userRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.UserRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/service.TokenPair"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    }
                }
            }
        },
        "/persons": {
            "get": {
                "description": "Get all persons",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Person"
                ],
                "summary": "Get all persons",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer token for authentication",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/model.Person"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    }
                }
            },
            "post": {
                "description": "Creates a new person",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Person"
                ],
                "summary": "Create a new person",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer token for authentication",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "personRequest value (model.PersonRequest)",
                        "name": "personRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.PersonRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/model.Person"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    }
                }
            }
        },
        "/persons/{id}": {
            "get": {
                "description": "Get a person by ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Person"
                ],
                "summary": "Get a person by ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer token for authentication",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Person ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Person"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    }
                }
            },
            "put": {
                "description": "Update a person by ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Person"
                ],
                "summary": "Update a person by ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer token for authentication",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Person ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "personRequest value (model.PersonRequest)",
                        "name": "personRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.PersonRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Person"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    }
                }
            },
            "delete": {
                "description": "Delete a person by ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Person"
                ],
                "summary": "Delete a person by ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer token for authentication",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Person ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    }
                }
            }
        },
        "/refresh": {
            "post": {
                "description": "Refreshes the access and refresh tokens using the provided refresh token",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Authentication"
                ],
                "summary": "Refreshes access and refresh tokens",
                "parameters": [
                    {
                        "description": "Refresh token",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/service.TokenPair"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    }
                }
            }
        },
        "/signUp": {
            "post": {
                "description": "Sign up a new user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Sign up a new user",
                "parameters": [
                    {
                        "description": "userRequest value (model.PersonRequest)",
                        "name": "userRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.UserRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    }
                }
            }
        }
    },
    "definitions": {
        "model.Person": {
            "type": "object",
            "required": [
                "profession",
                "salary"
            ],
            "properties": {
                "ID": {
                    "type": "string"
                },
                "married": {
                    "type": "boolean"
                },
                "profession": {
                    "type": "string",
                    "maxLength": 30,
                    "minLength": 3
                },
                "salary": {
                    "type": "integer",
                    "maximum": 100000,
                    "minimum": 100
                }
            }
        },
        "model.PersonRequest": {
            "type": "object",
            "required": [
                "profession",
                "salary"
            ],
            "properties": {
                "married": {
                    "type": "boolean"
                },
                "profession": {
                    "type": "string",
                    "maxLength": 30,
                    "minLength": 3
                },
                "salary": {
                    "type": "integer",
                    "maximum": 100000,
                    "minimum": 100
                }
            }
        },
        "model.UserRequest": {
            "type": "object",
            "required": [
                "password",
                "username"
            ],
            "properties": {
                "password": {
                    "type": "string",
                    "maxLength": 15,
                    "minLength": 4
                },
                "username": {
                    "type": "string",
                    "maxLength": 15,
                    "minLength": 4
                }
            }
        },
        "service.TokenPair": {
            "type": "object",
            "properties": {
                "accessToken": {
                    "type": "string"
                },
                "refreshToken": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "FirstTask API",
	Description:      "API for managing persons and users",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
