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
			"id": "03oa5ildtnmrxsn",
			"created": "2023-06-29 22:42:12.279Z",
			"updated": "2023-06-29 22:42:12.279Z",
			"name": "answers",
			"type": "base",
			"system": false,
			"schema": [
				{
					"system": false,
					"id": "c6jkot5j",
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
				},
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
		}`

		collection := &models.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collection); err != nil {
			return err
		}

		return daos.New(db).SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("03oa5ildtnmrxsn")
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	})
}
