package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tools/types"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("pcfz1rdnine760h")
		if err != nil {
			return err
		}

		collection.ListRule = nil

		collection.CreateRule = nil

		// remove
		collection.Schema.RemoveField("sjxisor6")

		// add
		new_game := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "gjd6k5tg",
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

		collection, err := dao.FindCollectionByNameOrId("pcfz1rdnine760h")
		if err != nil {
			return err
		}

		collection.ListRule = types.Pointer("@request.headers.x_pocketbase_game = game.code && game.puzzle = puzzle.id")

		collection.CreateRule = types.Pointer("@request.headers.x_pocketbase_game = game.code && game.puzzle = puzzle.id")

		// add
		del_game := &schema.SchemaField{}
		json.Unmarshal([]byte(`{
			"system": false,
			"id": "sjxisor6",
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
		collection.Schema.RemoveField("gjd6k5tg")

		return dao.SaveCollection(collection)
	})
}
