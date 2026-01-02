package storage

import (
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func newMongoStorage(mt *mtest.T) *MongoStorage {
	return &MongoStorage{
		collection: mt.Coll,
	}
}

func TestMongoStorage_Save(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		ms := newMongoStorage(mt)
		err := ms.Save("k1", []byte("value"))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestMongoStorage_Save_Error(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("replace error", func(mt *mtest.T) {
		mt.AddMockResponses(
			mtest.CreateCommandErrorResponse(
				mtest.CommandError{
					Code:    1,
					Message: "replace failed",
				},
			),
		)

		ms := newMongoStorage(mt)
		err := ms.Save("k", []byte("v"))
		if err == nil {
			t.Fatalf("expected error")
		}
	})
}

func TestMongoStorage_Load_Success(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("found", func(mt *mtest.T) {
		mt.AddMockResponses(
			mtest.CreateCursorResponse(
				1,
				"db.coll",
				mtest.FirstBatch,
				bson.D{
					{Key: "_id", Value: "k"},
					{Key: "data", Value: []byte("value")},
				},
			),
		)

		ms := newMongoStorage(mt)
		data, err := ms.Load("k")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if string(data) != "value" {
			t.Fatalf("unexpected data")
		}
	})
}

func TestMongoStorage_Load_NotFound(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("not found", func(mt *mtest.T) {
		mt.AddMockResponses(
			mtest.CreateCursorResponse(
				0,
				"db.coll",
				mtest.FirstBatch,
			),
		)

		ms := newMongoStorage(mt)
		_, err := ms.Load("missing")
		if err == nil {
			t.Fatalf("expected error")
		}
	})
}

func TestMongoStorage_Load_Error(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("find error", func(mt *mtest.T) {
		mt.AddMockResponses(
			mtest.CreateCommandErrorResponse(
				mtest.CommandError{
					Code:    1,
					Message: "find failed",
				},
			),
		)

		ms := newMongoStorage(mt)
		_, err := ms.Load("k")
		if err == nil {
			t.Fatalf("expected error")
		}
	})
}

func TestMongoStorage_Delete(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		ms := newMongoStorage(mt)
		err := ms.Delete("k")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestMongoStorage_Delete_Error(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("delete error", func(mt *mtest.T) {
		mt.AddMockResponses(
			mtest.CreateCommandErrorResponse(
				mtest.CommandError{
					Code:    1,
					Message: "delete failed",
				},
			),
		)

		ms := newMongoStorage(mt)
		err := ms.Delete("k")
		if err == nil {
			t.Fatalf("expected error")
		}
	})
}

func TestMongoStorage_Exists(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("exists true", func(mt *mtest.T) {
		mt.AddMockResponses(
			mtest.CreateCursorResponse(
				1,
				"db.coll",
				mtest.FirstBatch,
				bson.D{{Key: "n", Value: int64(1)}},
			),
		)

		ms := newMongoStorage(mt)
		ok, err := ms.Exists("k")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !ok {
			t.Fatalf("expected exists")
		}
	})
}

func TestMongoStorage_Exists_Error(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("count error", func(mt *mtest.T) {
		mt.AddMockResponses(
			mtest.CreateCommandErrorResponse(
				mtest.CommandError{
					Code:    1,
					Message: "count failed",
				},
			),
		)

		ms := newMongoStorage(mt)
		_, err := ms.Exists("k")
		if err == nil {
			t.Fatalf("expected error")
		}
	})
}

func TestMongoStorage_Load_MapsErrNoDocuments(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("mongo no documents", func(mt *mtest.T) {
		mt.AddMockResponses(
			mtest.CreateCursorResponse(
				0,
				"db.coll",
				mtest.FirstBatch,
			),
		)

		ms := newMongoStorage(mt)
		_, err := ms.Load("k")
		if err == nil {
			t.Fatalf("expected error")
		}
		if err.Error() != "key not found" {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
