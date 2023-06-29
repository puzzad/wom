package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		jsonData := `[
			{
				"id": "_pb_users_auth_",
				"created": "2023-06-26 19:27:36.595Z",
				"updated": "2023-06-29 21:19:29.367Z",
				"name": "users",
				"type": "auth",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "users_name",
						"name": "name",
						"type": "text",
						"required": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "users_avatar",
						"name": "avatar",
						"type": "file",
						"required": false,
						"unique": false,
						"options": {
							"maxSelect": 1,
							"maxSize": 5242880,
							"mimeTypes": [
								"image/jpeg",
								"image/png",
								"image/svg+xml",
								"image/gif",
								"image/webp"
							],
							"thumbs": null,
							"protected": false
						}
					},
					{
						"system": false,
						"id": "0rywdkyy",
						"name": "games",
						"type": "relation",
						"required": false,
						"unique": false,
						"options": {
							"collectionId": "vztgyvjzre4vxaf",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": null,
							"displayFields": []
						}
					}
				],
				"indexes": [],
				"listRule": "id = @request.auth.id",
				"viewRule": "id = @request.auth.id",
				"createRule": "",
				"updateRule": "id = @request.auth.id",
				"deleteRule": "id = @request.auth.id",
				"options": {
					"allowEmailAuth": true,
					"allowOAuth2Auth": true,
					"allowUsernameAuth": false,
					"exceptEmailDomains": null,
					"manageRule": null,
					"minPasswordLength": 8,
					"onlyEmailDomains": null,
					"requireEmail": true
				}
			},
			{
				"id": "8480lghxmlrhtn6",
				"created": "2023-06-26 19:53:11.369Z",
				"updated": "2023-06-29 20:11:18.444Z",
				"name": "hints",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "kqwwrg9z",
						"name": "title",
						"type": "text",
						"required": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "lxpw9omw",
						"name": "message",
						"type": "text",
						"required": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "vm40kxms",
						"name": "puzzle",
						"type": "relation",
						"required": false,
						"unique": false,
						"options": {
							"collectionId": "k5593ds7n07c487",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": []
						}
					},
					{
						"system": false,
						"id": "8fh5a7n7",
						"name": "order",
						"type": "number",
						"required": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null
						}
					}
				],
				"indexes": [],
				"listRule": null,
				"viewRule": null,
				"createRule": null,
				"updateRule": null,
				"deleteRule": null,
				"options": {}
			},
			{
				"id": "mlrnksy9yxm2dmv",
				"created": "2023-06-26 19:53:51.292Z",
				"updated": "2023-06-29 20:11:18.444Z",
				"name": "mailinglist",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "mxzs2jxw",
						"name": "email",
						"type": "text",
						"required": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "ivpsgwyd",
						"name": "subscribed",
						"type": "date",
						"required": false,
						"unique": false,
						"options": {
							"min": "",
							"max": ""
						}
					}
				],
				"indexes": [],
				"listRule": null,
				"viewRule": null,
				"createRule": null,
				"updateRule": null,
				"deleteRule": null,
				"options": {}
			},
			{
				"id": "t6dabzlosjg7sdb",
				"created": "2023-06-26 19:54:33.247Z",
				"updated": "2023-06-29 21:18:48.907Z",
				"name": "usedhints",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "l1n7g3rq",
						"name": "usedat",
						"type": "date",
						"required": false,
						"unique": false,
						"options": {
							"min": "",
							"max": ""
						}
					},
					{
						"system": false,
						"id": "yrjkj7b8",
						"name": "hint",
						"type": "relation",
						"required": false,
						"unique": false,
						"options": {
							"collectionId": "8480lghxmlrhtn6",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": []
						}
					},
					{
						"system": false,
						"id": "5cv8lsrq",
						"name": "game",
						"type": "relation",
						"required": false,
						"unique": false,
						"options": {
							"collectionId": "vztgyvjzre4vxaf",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": []
						}
					}
				],
				"indexes": [],
				"listRule": null,
				"viewRule": null,
				"createRule": null,
				"updateRule": null,
				"deleteRule": null,
				"options": {}
			},
			{
				"id": "dsyy96h6bthpiev",
				"created": "2023-06-26 19:55:28.099Z",
				"updated": "2023-06-29 20:11:18.444Z",
				"name": "adventures",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "bscbu4lg",
						"name": "name",
						"type": "text",
						"required": true,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "uisegeh8",
						"name": "description",
						"type": "editor",
						"required": true,
						"unique": false,
						"options": {}
					},
					{
						"system": false,
						"id": "qwkmayia",
						"name": "price",
						"type": "number",
						"required": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null
						}
					},
					{
						"system": false,
						"id": "jk6czf48",
						"name": "public",
						"type": "bool",
						"required": false,
						"unique": false,
						"options": {}
					},
					{
						"system": false,
						"id": "1hhjmvpi",
						"name": "firstpuzzle",
						"type": "relation",
						"required": false,
						"unique": false,
						"options": {
							"collectionId": "k5593ds7n07c487",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": []
						}
					},
					{
						"system": false,
						"id": "4j7cklhe",
						"name": "background",
						"type": "file",
						"required": false,
						"unique": false,
						"options": {
							"maxSelect": 1,
							"maxSize": 5242880,
							"mimeTypes": [],
							"thumbs": [],
							"protected": false
						}
					},
					{
						"system": false,
						"id": "enj9csmb",
						"name": "logo",
						"type": "file",
						"required": false,
						"unique": false,
						"options": {
							"maxSelect": 1,
							"maxSize": 5242880,
							"mimeTypes": [],
							"thumbs": [],
							"protected": false
						}
					},
					{
						"system": false,
						"id": "n9qydvmn",
						"name": "features",
						"type": "json",
						"required": false,
						"unique": false,
						"options": {}
					}
				],
				"indexes": [],
				"listRule": "public=true",
				"viewRule": "public=true",
				"createRule": null,
				"updateRule": null,
				"deleteRule": null,
				"options": {}
			},
			{
				"id": "pcfz1rdnine760h",
				"created": "2023-06-26 19:58:21.418Z",
				"updated": "2023-06-29 22:05:11.098Z",
				"name": "guesses",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "gjjmekcy",
						"name": "content",
						"type": "text",
						"required": true,
						"unique": false,
						"options": {
							"min": null,
							"max": 50,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "htkjhtym",
						"name": "puzzle",
						"type": "relation",
						"required": false,
						"unique": false,
						"options": {
							"collectionId": "k5593ds7n07c487",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": []
						}
					},
					{
						"system": false,
						"id": "hukqqnrx",
						"name": "correct",
						"type": "bool",
						"required": false,
						"unique": false,
						"options": {}
					},
					{
						"system": false,
						"id": "gjd6k5tg",
						"name": "game",
						"type": "relation",
						"required": false,
						"unique": false,
						"options": {
							"collectionId": "vztgyvjzre4vxaf",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": []
						}
					}
				],
				"indexes": [],
				"listRule": "@request.auth.id = game.id && @request.auth.puzzle.id = puzzle.id",
				"viewRule": "@request.auth.id = game.id && @request.auth.puzzle.id = puzzle.id",
				"createRule": "@request.auth.id = game.id && @request.auth.puzzle.id = puzzle.id",
				"updateRule": null,
				"deleteRule": null,
				"options": {}
			},
			{
				"id": "k5593ds7n07c487",
				"created": "2023-06-26 19:58:51.790Z",
				"updated": "2023-06-29 21:51:59.418Z",
				"name": "puzzles",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "nutgcd6e",
						"name": "title",
						"type": "text",
						"required": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "3ktxtegh",
						"name": "next",
						"type": "relation",
						"required": false,
						"unique": false,
						"options": {
							"collectionId": "k5593ds7n07c487",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": []
						}
					},
					{
						"system": false,
						"id": "oolgqnvg",
						"name": "adventure",
						"type": "relation",
						"required": false,
						"unique": false,
						"options": {
							"collectionId": "dsyy96h6bthpiev",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": []
						}
					},
					{
						"system": false,
						"id": "5pkihs4z",
						"name": "information",
						"type": "editor",
						"required": false,
						"unique": false,
						"options": {}
					},
					{
						"system": false,
						"id": "w5epy1nu",
						"name": "story",
						"type": "editor",
						"required": false,
						"unique": false,
						"options": {}
					},
					{
						"system": false,
						"id": "myloln6s",
						"name": "puzzle",
						"type": "editor",
						"required": false,
						"unique": false,
						"options": {}
					}
				],
				"indexes": [],
				"listRule": null,
				"viewRule": "@request.auth.puzzle.id = id",
				"createRule": null,
				"updateRule": null,
				"deleteRule": null,
				"options": {}
			},
			{
				"id": "vztgyvjzre4vxaf",
				"created": "2023-06-29 21:16:45.535Z",
				"updated": "2023-06-29 21:30:09.792Z",
				"name": "games",
				"type": "auth",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "48eqclsh",
						"name": "status",
						"type": "select",
						"required": false,
						"unique": false,
						"options": {
							"maxSelect": 1,
							"values": [
								"UNPAID",
								"PAID",
								"EXPIRED"
							]
						}
					},
					{
						"system": false,
						"id": "9ifgyjsw",
						"name": "purchaser",
						"type": "relation",
						"required": false,
						"unique": false,
						"options": {
							"collectionId": "_pb_users_auth_",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": []
						}
					},
					{
						"system": false,
						"id": "j5rdpozc",
						"name": "adventure",
						"type": "relation",
						"required": true,
						"unique": false,
						"options": {
							"collectionId": "dsyy96h6bthpiev",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": []
						}
					},
					{
						"system": false,
						"id": "9v3s2zwf",
						"name": "start",
						"type": "date",
						"required": false,
						"unique": false,
						"options": {
							"min": "",
							"max": ""
						}
					},
					{
						"system": false,
						"id": "vpvrshty",
						"name": "end",
						"type": "date",
						"required": false,
						"unique": false,
						"options": {
							"min": "",
							"max": ""
						}
					},
					{
						"system": false,
						"id": "kodqnhxo",
						"name": "puzzle",
						"type": "relation",
						"required": false,
						"unique": false,
						"options": {
							"collectionId": "k5593ds7n07c487",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": []
						}
					}
				],
				"indexes": [],
				"listRule": "id = @request.auth.id || @request.auth.games.id ?= id",
				"viewRule": "id = @request.auth.id",
				"createRule": null,
				"updateRule": null,
				"deleteRule": null,
				"options": {
					"allowEmailAuth": false,
					"allowOAuth2Auth": false,
					"allowUsernameAuth": true,
					"exceptEmailDomains": [],
					"manageRule": null,
					"minPasswordLength": 5,
					"onlyEmailDomains": [],
					"requireEmail": false
				}
			}
		]`

		collections := []*models.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collections); err != nil {
			return err
		}

		return daos.New(db).ImportCollections(collections, true, nil)
	}, func(db dbx.Builder) error {
		return nil
	})
}
