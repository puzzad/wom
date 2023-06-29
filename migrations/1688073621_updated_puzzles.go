package migrations

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/tools/types"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("k5593ds7n07c487")
		if err != nil {
			return err
		}

		collection.ViewRule = nil

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("k5593ds7n07c487")
		if err != nil {
			return err
		}

		collection.ViewRule = types.Pointer("@collection.games.adventure.id ?= adventure.id  && @collection.games.puzzle.id ?= id && @collection.games.code ?= @request.headers.x_pocketbase_game")

		return dao.SaveCollection(collection)
	})
}
