// Package tape Code generated by swaggo/swag. DO NOT EDIT
package tape

import "github.com/swaggo/swag"

const docTemplatetape = `{
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
        "/eject": {
            "post": {
                "description": "eject tape",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tm"
                ],
                "summary": "eject tape",
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/erase": {
            "delete": {
                "description": "erase tape",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "tm"
                ],
                "summary": "erase tape",
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/info": {
            "get": {
                "description": "get tape info",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "info"
                ],
                "summary": "get tape info",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.TapeInfo"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/read": {
            "get": {
                "description": "read file",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "file"
                ],
                "summary": "read file",
                "parameters": [
                    {
                        "type": "array",
                        "items": {
                            "type": "integer"
                        },
                        "collectionFormat": "csv",
                        "description": "file numbers",
                        "name": "numbers",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "path to extract dir",
                        "name": "path",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/write": {
            "post": {
                "description": "write file",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "file"
                ],
                "summary": "write file",
                "parameters": [
                    {
                        "type": "string",
                        "description": "путь до файла",
                        "name": "filePath",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.FileWriteInfo"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        }
    },
    "definitions": {
        "model.Attributes": {
            "type": "object",
            "properties": {
                "assigningOrg": {
                    "type": "string"
                },
                "cartridgeLoadCount": {
                    "type": "integer"
                },
                "cartridgeType": {
                    "type": "string"
                },
                "mAMCapacity": {
                    "$ref": "#/definitions/model.MAMCapacityAttribute"
                },
                "manufactureDate": {
                    "type": "string"
                },
                "manufacturer": {
                    "type": "string"
                },
                "mediumDensity": {
                    "$ref": "#/definitions/model.MediumDensityAttribute"
                },
                "partCapMax": {
                    "description": "Всего места в байтах",
                    "type": "integer"
                },
                "partCapRemain": {
                    "description": "Свободное место в байтах",
                    "type": "integer"
                },
                "readSession": {
                    "type": "integer"
                },
                "serialNumber": {
                    "type": "string"
                },
                "sessions": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.SessionAttribute"
                    }
                },
                "specs": {
                    "$ref": "#/definitions/model.SpecsAttribute"
                },
                "tapeLength": {
                    "description": "Длинна ленты в метрах",
                    "type": "integer"
                },
                "tapeWidth": {
                    "description": "Ширина ленты в милиметрах",
                    "type": "integer"
                },
                "totalRead": {
                    "type": "integer"
                },
                "totalWritten": {
                    "type": "integer"
                },
                "writtenSession": {
                    "type": "integer"
                }
            }
        },
        "model.FileWriteInfo": {
            "type": "object",
            "properties": {
                "bytesWrite": {
                    "type": "integer"
                },
                "fileNumber": {
                    "type": "integer"
                }
            }
        },
        "model.MAMCapacityAttribute": {
            "type": "object",
            "properties": {
                "capacity": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "spaceRemaining": {
                    "type": "integer"
                }
            }
        },
        "model.MediumDensityAttribute": {
            "type": "object",
            "properties": {
                "formattedAs": {
                    "type": "string"
                },
                "mediumformat": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "model.SessionAttribute": {
            "type": "object",
            "properties": {
                "devname": {
                    "type": "string"
                },
                "number": {
                    "type": "integer"
                },
                "serial": {
                    "type": "string"
                }
            }
        },
        "model.SpecsAttribute": {
            "type": "object",
            "properties": {
                "capacity": {
                    "$ref": "#/definitions/model.SpecsCapacityAttribute"
                },
                "duration": {
                    "$ref": "#/definitions/model.SpecsDurationAttribute"
                },
                "partitions": {
                    "$ref": "#/definitions/model.SpecsPartitionsAttribute"
                },
                "phy": {
                    "$ref": "#/definitions/model.SpecsPhyAttribute"
                },
                "speed": {
                    "$ref": "#/definitions/model.SpecsSpeedAttribute"
                }
            }
        },
        "model.SpecsCapacityAttribute": {
            "type": "object",
            "properties": {
                "compressFactor": {
                    "type": "string"
                },
                "compressed": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "native": {
                    "type": "integer"
                }
            }
        },
        "model.SpecsDurationAttribute": {
            "type": "object",
            "properties": {
                "fullTapeMinutes": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "model.SpecsPartitionsAttribute": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                },
                "partitionNumber": {
                    "type": "integer"
                }
            }
        },
        "model.SpecsPhyAttribute": {
            "type": "object",
            "properties": {
                "bandsPerTape": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "total": {
                    "type": "integer"
                },
                "tracksPerWrap": {
                    "type": "integer"
                },
                "wrapsPerBand": {
                    "type": "integer"
                }
            }
        },
        "model.SpecsSpeedAttribute": {
            "type": "object",
            "properties": {
                "compressed": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "native": {
                    "type": "integer"
                }
            }
        },
        "model.TapeInfo": {
            "type": "object",
            "properties": {
                "attributes": {
                    "$ref": "#/definitions/model.Attributes"
                },
                "firmware": {
                    "type": "string"
                },
                "model": {
                    "type": "string"
                },
                "vendor": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfotape holds exported Swagger Info so clients can modify it
var SwaggerInfotape = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Tape api",
	Description:      "Tape manage server.",
	InfoInstanceName: "tape",
	SwaggerTemplate:  docTemplatetape,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfotape.InstanceName(), SwaggerInfotape)
}
