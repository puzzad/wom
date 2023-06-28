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

		collection, err := dao.FindCollectionByNameOrId("_pb_users_auth_")
		if err != nil {
			return err
		}

		// add
		new_games := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "nqcdf6fn",
			"name": "games",
			"type": "relation",
			"required": false,
			"unique": false,
			"options": {
				"collectionId": "410a6wyrxhi89cg",
				"cascadeDelete": false,
				"minSelect": null,
				"maxSelect": null,
				"displayFields": []
			}
		}`), new_games)
		collection.Schema.AddField(new_games)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("_pb_users_auth_")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("nqcdf6fn")

		return dao.SaveCollection(collection)
	})
}
