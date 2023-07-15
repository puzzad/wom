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
		dao := daos.New(db)

		collection, err := dao.FindCollectionByNameOrId("djlovsu233v2617")
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	}, func(db dbx.Builder) error {
		jsonData := `{
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
			}`

		collection := &models.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collection); err != nil {
			return err
		}

		return daos.New(db).SaveCollection(collection)
	})
}
