// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag at
// 2019-04-12 16:03:13.187928 +0900 KST m=+0.069815385

package docs

import (
	"bytes"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": [
        "https"
    ],
    "swagger": "2.0",
    "info": {
        "description": "This is a Sumwhere server API",
        "title": "Sumwhere API",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "url": "https://www.sumwhere.kr",
            "email": "qjadn0914@naver.com"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "2.0"
    },
    "host": "www.sumwhere.kr",
    "basePath": "/v1",
    "paths": {
        "/admin/assgin": {
            "patch": {
                "description": "기존 어드민사용자가 추가로 어드민 권한을 유저에게 ≠부여",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "admin"
                ],
                "summary": "어드민 권한 부여",
                "parameters": [
                    {
                        "description": "{'id':1}",
                        "name": "id",
                        "in": "body",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/controllers.AdminController"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "해당 유저정보",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/models.User"
                        }
                    }
                }
            }
        },
        "/signin/email": {
            "post": {
                "description": "이메일을 이용한 로그인",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "signin"
                ],
                "summary": "이메일 로그인",
                "parameters": [
                    {
                        "description": "email 과 passward만 쓸것",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/models.User"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "access 토큰을 반환",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/controllers.Token"
                        }
                    }
                }
            }
        },
        "/signin/facebook": {
            "post": {
                "description": "페이스북을 이용한 로그인",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "signin"
                ],
                "summary": "페이스북 로그인",
                "parameters": [
                    {
                        "description": "토큰 전송 ",
                        "name": "token",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/controllers.Token"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "access 토큰을 반환",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/controllers.Token"
                        }
                    }
                }
            }
        },
        "/signin/kakao": {
            "post": {
                "description": "카카오를 이용한 로그인",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "signin"
                ],
                "summary": "kakao 로그인",
                "parameters": [
                    {
                        "description": "토큰 전송 ",
                        "name": "token",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/controllers.Token"
                        }
                    },
                    {
                        "description": "Bottle ID",
                        "name": "token",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/controllers.Token"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "access 토큰을 반환",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/controllers.Token"
                        }
                    }
                }
            }
        },
        "/signup/email": {
            "post": {
                "description": "이메일을 이용한 가입",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "signup"
                ],
                "summary": "이메일 가입",
                "parameters": [
                    {
                        "description": "email 과 passward만 쓸것",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/models.User"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "access 토큰을 반환",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/controllers.Token"
                        }
                    }
                }
            }
        },
        "/signup/nickname": {
            "get": {
                "description": "닉네임이 존재하는지 확인",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "signup"
                ],
                "summary": "닉네임 확인",
                "responses": {
                    "200": {
                        "description": "bool 타입의 결과 있으면 true 없으면 false",
                        "schema": {
                            "type": "boolean"
                        }
                    }
                }
            }
        },
        "/user/all": {
            "get": {
                "description": "가입자 리스트 반환",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "유저 리스트 반환",
                "parameters": [
                    {
                        "description": "순서,순서마다의 정렬순서 (desc|asc), 어디서부터, 몇개를 가져올지",
                        "name": "query",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/models.GetQuery"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "쿼리의 결과",
                        "schema": {
                            "type": "object",
                            "$ref": "#/definitions/utils.ArrayResult"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "controllers.AdminController": {
            "type": "object"
        },
        "controllers.Token": {
            "type": "object",
            "properties": {
                "access_token": {
                    "type": "string",
                    "example": "ldifgj1lij31t9gsegl"
                }
            }
        },
        "models.GetQuery": {
            "type": "object",
            "properties": {
                "limit": {
                    "type": "integer",
                    "example": 2
                },
                "offset": {
                    "type": "integer",
                    "example": 1
                },
                "orderBy": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    },
                    "example": [
                        "desc"
                    ]
                },
                "sortBy": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    },
                    "example": [
                        "email"
                    ]
                }
            }
        },
        "models.User": {
            "type": "object",
            "properties": {
                "age": {
                    "type": "integer"
                },
                "createdAt": {
                    "type": "string"
                },
                "deletedAt": {
                    "type": "string"
                },
                "email": {
                    "type": "string",
                    "example": "example@naver.com"
                },
                "gender": {
                    "type": "string"
                },
                "hasProfile": {
                    "type": "boolean"
                },
                "id": {
                    "type": "integer"
                },
                "is_admin": {
                    "type": "boolean"
                },
                "joinType": {
                    "type": "string"
                },
                "mainProfileImage": {
                    "type": "string"
                },
                "nickname": {
                    "type": "string"
                },
                "password": {
                    "type": "string",
                    "example": "1234qwer"
                },
                "point": {
                    "type": "integer"
                },
                "snsId": {
                    "type": "string"
                },
                "token": {
                    "type": "string"
                },
                "updatedAt": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "utils.ArrayResult": {
            "type": "object",
            "properties": {
                "items": {
                    "type": "object"
                },
                "totalCount": {
                    "type": "integer"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo swaggerInfo

type s struct{}

func (s *s) ReadDoc() string {
	t, err := template.New("swagger_info").Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, SwaggerInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
