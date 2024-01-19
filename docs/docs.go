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
        "/buy-store": {
            "post": {
                "description": "Buy the assests from the shop",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Store"
                ],
                "summary": "Buy things",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Player Access token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "shop Details",
                        "name": "loginDetails",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.BuyStoreRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Sucess",
                        "schema": {
                            "$ref": "#/definitions/response.Success"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "$ref": "#/definitions/response.Success"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/response.Success"
                        }
                    }
                }
            }
        },
        "/check-otp": {
            "post": {
                "description": "Verifies the otp sent on email",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Authentication"
                ],
                "summary": "Verifies OTP",
                "parameters": [
                    {
                        "description": "Email Details",
                        "name": "loginDetails",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.OtpRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Login successful",
                        "schema": {
                            "$ref": "#/definitions/response.Success"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "$ref": "#/definitions/response.Success"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/response.Success"
                        }
                    }
                }
            }
        },
        "/get-daily-goals": {
            "get": {
                "description": "Gets the daily goals for given player",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "DailyGoal"
                ],
                "summary": "Get the daily goals",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Player Access token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Sucess",
                        "schema": {
                            "$ref": "#/definitions/response.Success"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "$ref": "#/definitions/response.Success"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/response.Success"
                        }
                    }
                }
            }
        },
        "/get-level-rewards": {
            "get": {
                "description": "Gets the rewards according to player level",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Player"
                ],
                "summary": "Gets reward list",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Player Access token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Login successful",
                        "schema": {
                            "$ref": "#/definitions/response.Success"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "$ref": "#/definitions/response.Success"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/response.Success"
                        }
                    }
                }
            }
        },
        "/get-settings": {
            "get": {
                "description": "Gets the current settings of that player",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Player"
                ],
                "summary": "Gets the settings",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Player Access token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Sucess",
                        "schema": {
                            "$ref": "#/definitions/response.Success"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "$ref": "#/definitions/response.Success"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/response.Success"
                        }
                    }
                }
            }
        },
        "/level-reward-collect": {
            "post": {
                "description": "Collects the reward for a level of that player",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Player"
                ],
                "summary": "Collects Reward",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Player Access token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "Player Details",
                        "name": "loginDetails",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.PlayerLevelRewardCollectRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Login successful",
                        "schema": {
                            "$ref": "#/definitions/response.Success"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "$ref": "#/definitions/response.Success"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/response.Success"
                        }
                    }
                }
            }
        },
        "/popupoffers": {
            "get": {
                "description": "Get the specific type of reward",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Store"
                ],
                "summary": "Get the specific type of reward",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Type of reward",
                        "name": "type",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success",
                        "schema": {
                            "$ref": "#/definitions/response.Success"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "$ref": "#/definitions/response.Success"
                        }
                    },
                    "401": {
                        "description": "Unauthorised",
                        "schema": {
                            "$ref": "#/definitions/response.Success"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/response.Success"
                        }
                    }
                }
            }
        },
        "/reset-password": {
            "post": {
                "description": "Resets the password of the user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Authentication"
                ],
                "summary": "Resets the password",
                "parameters": [
                    {
                        "description": "Email Details",
                        "name": "loginDetails",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.RestPasswordRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Login successful",
                        "schema": {
                            "$ref": "#/definitions/response.Success"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "$ref": "#/definitions/response.Success"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/response.Success"
                        }
                    }
                }
            }
        },
        "/send-otp": {
            "post": {
                "description": "Sends the otp on the register email",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Authentication"
                ],
                "summary": "Sends OTP",
                "parameters": [
                    {
                        "description": "Email Details",
                        "name": "loginDetails",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.EmailRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Login successful",
                        "schema": {
                            "$ref": "#/definitions/response.Success"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "$ref": "#/definitions/response.Success"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/response.Success"
                        }
                    }
                }
            }
        },
        "/stats": {
            "get": {
                "description": "Get the player game stats",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Player"
                ],
                "summary": "Get player game stats",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Player Access token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Login successful",
                        "schema": {
                            "$ref": "#/definitions/response.Success"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "$ref": "#/definitions/response.Success"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/response.Success"
                        }
                    }
                }
            }
        },
        "/store": {
            "get": {
                "description": "Gets shop details",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Store"
                ],
                "summary": "Gets shop details",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Player Access token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Sucess",
                        "schema": {
                            "$ref": "#/definitions/response.Success"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "$ref": "#/definitions/response.Success"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/response.Success"
                        }
                    }
                }
            }
        },
        "/update-settings": {
            "put": {
                "description": "Updates the game settings of that player",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Player"
                ],
                "summary": "Updates setting",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Player Access token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "Player Details",
                        "name": "loginDetails",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.PlayerLevelRewardCollectRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Sucess",
                        "schema": {
                            "$ref": "#/definitions/response.Success"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "$ref": "#/definitions/response.Success"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/response.Success"
                        }
                    }
                }
            }
        },
        "/user-data": {
            "get": {
                "description": "Get game stats of ther players",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Player"
                ],
                "summary": "Get game stats of ther players",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Id of the other player",
                        "name": "playerId",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Login successful",
                        "schema": {
                            "$ref": "#/definitions/response.Success"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "$ref": "#/definitions/response.Success"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/response.Success"
                        }
                    }
                }
            },
            "put": {
                "description": "Updates player info like username and avatar",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Player"
                ],
                "summary": "Updates player info",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Player Access token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "Player Details",
                        "name": "loginDetails",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.UpdatePlayer"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Login successful",
                        "schema": {
                            "$ref": "#/definitions/response.Success"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "$ref": "#/definitions/response.Success"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/response.Success"
                        }
                    }
                }
            }
        },
        "/users/email-verify": {
            "get": {
                "description": "Perform Email verifictaion",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Authentication"
                ],
                "summary": "User email verification",
                "responses": {
                    "200": {
                        "description": "Login successful",
                        "schema": {
                            "$ref": "#/definitions/response.Success"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "$ref": "#/definitions/response.Success"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/response.Success"
                        }
                    }
                }
            }
        },
        "/users/sign-in": {
            "post": {
                "description": "Perform Users login and generate access token",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Authentication"
                ],
                "summary": "User Login",
                "parameters": [
                    {
                        "description": "Login Details",
                        "name": "loginDetails",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.LoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Login successful",
                        "schema": {
                            "$ref": "#/definitions/response.Success"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "$ref": "#/definitions/response.Success"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/response.Success"
                        }
                    }
                }
            }
        },
        "/users/sign-out": {
            "delete": {
                "description": "Logs out a player",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Authentication"
                ],
                "summary": "Logout Player",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Player Access Token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/response.Success"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/response.Success"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/response.Success"
                        }
                    }
                }
            }
        },
        "/users/sign-up": {
            "post": {
                "description": "Perform signup and sends email for verification",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Authentication"
                ],
                "summary": "Sign Up",
                "parameters": [
                    {
                        "description": "Signup Request",
                        "name": "guestLoginRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.SigupRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Login successful",
                        "schema": {
                            "$ref": "#/definitions/response.Success"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "$ref": "#/definitions/response.Success"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/response.Success"
                        }
                    }
                }
            }
        },
        "/users/social-login": {
            "post": {
                "description": "Perform Users social login and generate access token",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Authentication"
                ],
                "summary": "User Login",
                "parameters": [
                    {
                        "description": "Login Details",
                        "name": "loginDetails",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/request.SocialLoginReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Login successful",
                        "schema": {
                            "$ref": "#/definitions/response.Success"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "$ref": "#/definitions/response.Success"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/response.Success"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "request.BuyStoreRequest": {
            "type": "object",
            "properties": {
                "productId": {
                    "type": "string"
                },
                "token": {
                    "type": "string"
                }
            }
        },
        "request.EmailRequest": {
            "type": "object",
            "properties": {
                "user": {
                    "type": "object",
                    "required": [
                        "email"
                    ],
                    "properties": {
                        "email": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "request.LoginRequest": {
            "type": "object",
            "properties": {
                "user": {
                    "type": "object",
                    "properties": {
                        "credential": {
                            "type": "string"
                        },
                        "password": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "request.OtpRequest": {
            "type": "object",
            "required": [
                "email",
                "otp"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "otp": {
                    "type": "integer"
                }
            }
        },
        "request.PlayerLevelRewardCollectRequest": {
            "type": "object",
            "properties": {
                "level": {
                    "type": "integer"
                }
            }
        },
        "request.RestPasswordRequest": {
            "type": "object",
            "required": [
                "email",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "request.SigupRequest": {
            "type": "object",
            "properties": {
                "user": {
                    "type": "object",
                    "properties": {
                        "avatar": {
                            "type": "integer"
                        },
                        "email": {
                            "type": "string"
                        },
                        "password": {
                            "type": "string"
                        },
                        "username": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "request.SocialLoginReq": {
            "type": "object",
            "properties": {
                "avatar": {
                    "type": "integer"
                },
                "email": {
                    "type": "string"
                },
                "uid": {
                    "type": "string"
                }
            }
        },
        "request.UpdatePlayer": {
            "type": "object",
            "properties": {
                "avatar": {
                    "type": "integer"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "response.Status": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "data": {},
                "message": {
                    "type": "string"
                }
            }
        },
        "response.Success": {
            "type": "object",
            "properties": {
                "status": {
                    "$ref": "#/definitions/response.Status"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/api/v1/",
	Schemes:          []string{},
	Title:            "Survival Party",
	Description:      "This is the api documentation of survival party game",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
