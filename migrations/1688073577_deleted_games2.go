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
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("410a6wyrxhi89cg")
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	}, func(db dbx.Builder) error {
		jsonData := `{
			"id": "410a6wyrxhi89cg",
			"created": "2023-06-26 19:57:50.207Z",
			"updated": "2023-06-29 21:19:04.532Z",
			"name": "games2",
			"type": "base",
			"system": false,
			"schema": [
				{
					"system": false,
					"id": "qigcpdpr",
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
					"id": "a4nt925u",
					"name": "code",
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
					"id": "h8yfzu0u",
					"name": "user",
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
					"id": "ooicceyv",
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
					"id": "ejdrirte",
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
					"id": "l50ob5lu",
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
					"id": "4wxfchd9",
					"name": "end",
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
		}`

		collection := &models.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collection); err != nil {
			return err
		}

		return daos.New(db).SaveCollection(collection)
	})
}
