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

		collection, err := dao.FindCollectionByNameOrId("vztgyvjzre4vxaf")
		if err != nil {
			return err
		}

		// add
		new_puzzle := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
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
		}`), new_puzzle)
		collection.Schema.AddField(new_puzzle)

		// update
		edit_adventure := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
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
		}`), edit_adventure)
		collection.Schema.AddField(edit_adventure)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("vztgyvjzre4vxaf")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("kodqnhxo")

		// update
		edit_adventure := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
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
		}`), edit_adventure)
		collection.Schema.AddField(edit_adventure)

		return dao.SaveCollection(collection)
	})
}
