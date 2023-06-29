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

		collection, err := dao.FindCollectionByNameOrId("bh86snmtso4vuoo")
		if err != nil {
			return err
		}

		// update
		edit_answer := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "j7q9yqpv",
			"name": "answer",
			"type": "text",
			"required": false,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"pattern": ""
			}
		}`), edit_answer)
		collection.Schema.AddField(edit_answer)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("bh86snmtso4vuoo")
		if err != nil {
			return err
		}

		// update
		edit_answer := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
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
		}`), edit_answer)
		collection.Schema.AddField(edit_answer)

		return dao.SaveCollection(collection)
	})
}
