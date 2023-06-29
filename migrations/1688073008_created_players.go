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
			"id": "g7uljy6ajppp2hb",
			"created": "2023-06-29 21:10:08.636Z",
			"updated": "2023-06-29 21:10:08.636Z",
			"name": "players",
			"type": "auth",
			"system": false,
			"schema": [
				{
					"system": false,
					"id": "ihdiomc2",
					"name": "game",
					"type": "relation",
					"required": false,
					"unique": false,
					"options": {
						"collectionId": "410a6wyrxhi89cg",
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
			"options": {
				"allowEmailAuth": false,
				"allowOAuth2Auth": false,
				"allowUsernameAuth": true,
				"exceptEmailDomains": [],
				"manageRule": null,
				"minPasswordLength": 8,
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

		collection, err := dao.FindCollectionByNameOrId("g7uljy6ajppp2hb")
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	})
}
