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

		collection, err := dao.FindCollectionByNameOrId("t6dabzlosjg7sdb")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("mclraxps")

		// add
		new_game := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
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
		}`), new_game)
		collection.Schema.AddField(new_game)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("t6dabzlosjg7sdb")
		if err != nil {
			return err
		}

		// add
		del_game := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "mclraxps",
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
		}`), del_game)
		collection.Schema.AddField(del_game)

		// remove
		collection.Schema.RemoveField("5cv8lsrq")

		return dao.SaveCollection(collection)
	})
}
