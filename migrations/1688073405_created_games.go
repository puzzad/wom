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
		jsonData := `{
			"id": "vztgyvjzre4vxaf",
			"created": "2023-06-29 21:16:45.535Z",
			"updated": "2023-06-29 21:16:45.535Z",
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
				}
			],
			"indexes": [],
			"listRule": null,
			"viewRule": null,
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
		}`

		collection := &models.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collection); err != nil {
			return err
		}

		return daos.New(db).SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("vztgyvjzre4vxaf")
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	})
}
