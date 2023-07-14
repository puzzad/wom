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
				"id": "8480lghxmlrhtn6",
				"created": "2023-07-08 16:30:10.635Z",
				"updated": "2023-07-10 18:36:37.726Z",
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
				"created": "2023-07-08 16:30:10.635Z",
				"updated": "2023-07-10 18:36:37.726Z",
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
				"created": "2023-07-08 16:30:10.636Z",
				"updated": "2023-07-10 18:36:37.726Z",
				"name": "usedhints",
				"type": "base",
				"system": false,
				"schema": [
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
				"created": "2023-07-08 16:30:10.636Z",
				"updated": "2023-07-10 18:36:37.726Z",
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
				"created": "2023-07-08 16:30:10.636Z",
				"updated": "2023-07-10 18:36:37.726Z",
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
				"created": "2023-07-08 16:30:10.636Z",
				"updated": "2023-07-10 18:36:37.726Z",
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
				"created": "2023-07-08 16:30:10.636Z",
				"updated": "2023-07-10 18:36:37.726Z",
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
			},
			{
				"id": "03oa5ildtnmrxsn",
				"created": "2023-07-08 16:30:10.637Z",
				"updated": "2023-07-10 18:36:37.727Z",
				"name": "answers",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "jyt3mvz0",
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
						"id": "hwuhbtvk",
						"name": "content",
						"type": "text",
						"required": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
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
				"id": "djlovsu233v2617",
				"created": "2023-07-08 23:24:57.624Z",
				"updated": "2023-07-10 18:36:37.728Z",
				"name": "currenthints",
				"type": "view",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "mtesco6l",
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
						"id": "mmmvgvfq",
						"name": "message",
						"type": "json",
						"required": false,
						"unique": false,
						"options": {}
					},
					{
						"system": false,
						"id": "cnrcq1kw",
						"name": "locked",
						"type": "json",
						"required": false,
						"unique": false,
						"options": {}
					},
					{
						"system": false,
						"id": "upzibpqs",
						"name": "puzzleid",
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
				"listRule": "@request.auth.puzzle.id = puzzleid",
				"viewRule": null,
				"createRule": null,
				"updateRule": null,
				"deleteRule": null,
				"options": {
					"query": "select hints.id as id, hints.title as title, iif(usedhints.id ISNULL, '', hints.message) as message, iif(usedhints.id ISNULL, true, false) as locked, hints.puzzle as puzzleid\nfrom hints\nleft join usedhints on hints.id = usedhints.hint"
				}
			},
			{
				"id": "_pb_users_auth_",
				"created": "2023-07-09 00:20:19.016Z",
				"updated": "2023-07-10 18:36:37.727Z",
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
				"id": "54fmr4z92ksrbwy",
				"created": "2023-07-09 19:06:47.721Z",
				"updated": "2023-07-10 18:36:37.729Z",
				"name": "solvetimes",
				"type": "view",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "ljgr6hhp",
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
						"id": "ot0xo1hg",
						"name": "timeSolved",
						"type": "json",
						"required": false,
						"unique": false,
						"options": {}
					},
					{
						"system": false,
						"id": "gltcykwz",
						"name": "gameStart",
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
						"id": "rqpkegbd",
						"name": "gameCode",
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
						"id": "4kilszrr",
						"name": "usedhints",
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
				"listRule": "@request.auth.username = gameCode",
				"viewRule": null,
				"createRule": null,
				"updateRule": null,
				"deleteRule": null,
				"options": {
					"query": "select guesses.id, puzzles.title, max(guesses.created) as timeSolved, games.start as gameStart, games.username as gameCode, count(usedhints.id) as usedhints\nfrom guesses\n         left join puzzles on guesses.puzzle = puzzles.id\n         left join games on guesses.game = games.id\n         left join hints on hints.puzzle = puzzles.id\n         left join usedhints on games.id = usedhints.game AND usedhints.game = games.id and usedhints.hint = hints.id\nwhere correct=true\nGROUP BY puzzles.id, games.username\norder by timeSolved"
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
