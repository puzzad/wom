package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models/schema"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("pcfz1rdnine760h")
		if err != nil {
			return err
		}

		// update
		edit_content := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
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
		}`), edit_content)
		collection.Schema.AddField(edit_content)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("pcfz1rdnine760h")
		if err != nil {
			return err
		}

		// update
		edit_content := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "gjjmekcy",
			"name": "content",
			"type": "text",
			"required": false,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"pattern": ""
			}
		}`), edit_content)
		collection.Schema.AddField(edit_content)

		return dao.SaveCollection(collection)
	})
}
