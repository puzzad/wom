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
			"id": "bh86snmtso4vuoo",
			"created": "2023-06-29 20:27:33.996Z",
			"updated": "2023-06-29 20:27:33.996Z",
			"name": "answers",
			"type": "base",
			"system": false,
			"schema": [
				{
					"system": false,
					"id": "zgyq0uxx",
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
					"id": "u0p5lha6",
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
				},
				{
					"system": false,
					"id": "j7q9yqpv",
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
		}`

		collection := &models.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collection); err != nil {
			return err
		}

		return daos.New(db).SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("bh86snmtso4vuoo")
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	})
}
