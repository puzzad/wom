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

		collection, err := dao.FindCollectionByNameOrId("pcfz1rdnine760h")
		if err != nil {
			return err
		}

		collection.ListRule = types.Pointer("@request.auth.id = game.id && @request.auth.puzzle.id = puzzle.id")

		collection.ViewRule = nil

		collection.CreateRule = types.Pointer("@request.auth.id = game.id && @request.auth.puzzle.id = puzzle.id")

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("pcfz1rdnine760h")
		if err != nil {
			return err
		}

		collection.ListRule = nil

		collection.ViewRule = types.Pointer("@request.auth.id = game.id && @request.auth.puzzle = puzzle.id")

		collection.CreateRule = nil

		return dao.SaveCollection(collection)
	})
}
