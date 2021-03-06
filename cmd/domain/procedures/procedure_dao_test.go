package procedures

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bxcodec/faker/v3"
	"github.com/ko1ke/know-sync-api/cmd/utils/rand_utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ProcedureDaoTestSuite struct {
	suite.Suite
	TestDB     *gorm.DB
	mock       sqlmock.Sqlmock
	dummy      Procedure
	dummyPItem ProcedureItem
}

// set up test to gorm send queries to mock
func (suite *ProcedureDaoTestSuite) SetupTest() {
	db, mock, _ := sqlmock.New()
	suite.mock = mock
	suite.TestDB, _ = gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
}

func (suite *ProcedureDaoTestSuite) TearDownTest() {
	db, _ := suite.TestDB.DB()
	db.Close()
}

func (suite *ProcedureDaoTestSuite) AfterTest(_, _ string) {
	require.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

// Run test suite
func TestProcedureDaoTestSuite(t *testing.T) {
	suite.Run(t, new(ProcedureDaoTestSuite))
}

func (suite *ProcedureDaoTestSuite) BeforeTest(suiteName string, testName string) {

	suite.dummy = Procedure{
		ID:              rand_utils.MakeRandomUInt(100),
		Title:           faker.Word(),
		Content:         faker.Sentence(),
		Publish:         true,
		UserID:          rand_utils.MakeRandomUInt(100),
		EyeCatchImgName: faker.Word() + ".jpeg",
	}
	// spew.Dump(suite.dummy)
}

func (suite *ProcedureDaoTestSuite) TestGet() {
	suite.Run("get a procedure", func() {
		suite.mock.ExpectQuery(
			regexp.QuoteMeta(`SELECT * FROM "procedures" WHERE id = $1 AND "procedures"."deleted_at" IS NULL AND "procedures"."id" = $2 ORDER BY "procedures"."id" LIMIT 1`),
		).WithArgs(suite.dummy.ID, suite.dummy.ID).
			WillReturnRows(suite.mock.NewRows([]string{"id", "title", "content", "user_id", "publish", "eye_catch_img_name"}).
				AddRow(suite.dummy.ID, suite.dummy.Title, suite.dummy.Content, suite.dummy.UserID, suite.dummy.Publish, suite.dummy.EyeCatchImgName))

		procedure := &Procedure{
			ID: suite.dummy.ID,
		}
		err := procedure.Get(suite.TestDB)
		require.NoError(suite.T(), err)
		assert.Equal(suite.T(), procedure.ID, suite.dummy.ID, "unexpected ID")
		assert.Equal(suite.T(), procedure.Title, suite.dummy.Title, "unexpected Title")
		assert.Equal(suite.T(), procedure.Content, suite.dummy.Content, "unexpected Content")
		assert.Equal(suite.T(), procedure.UserID, suite.dummy.UserID, "unexpected UserID")
		assert.Equal(suite.T(), procedure.Publish, suite.dummy.Publish, "unexpected Publish")
		assert.Equal(suite.T(), procedure.EyeCatchImgName, suite.dummy.EyeCatchImgName, "unexpected EyeCatchImgName")
	})
}

func (suite *ProcedureDaoTestSuite) TestGetItem() {
	suite.Run("get a procedure item", func() {
		suite.mock.ExpectQuery(
			regexp.QuoteMeta(`SELECT procedures.*, users.username FROM "procedures" inner join users ON users.id=procedures.user_id WHERE procedures.id = $1 AND "procedures"."deleted_at" IS NULL ORDER BY "procedures"."id" LIMIT 1`),
		).WithArgs(suite.dummy.ID).
			WillReturnRows(suite.mock.NewRows([]string{"id", "title", "content", "user_id", "users.username", "publish", "eye_catch_img_name"}).
				AddRow(suite.dummy.ID, suite.dummy.Title, suite.dummy.Content, suite.dummy.UserID, "", suite.dummy.Publish, suite.dummy.EyeCatchImgName))

		p, err := GetItem(suite.TestDB, suite.dummy.ID)

		require.NoError(suite.T(), err)
		assert.Equal(suite.T(), p.ID, suite.dummy.ID, "unexpected ID")
		assert.Equal(suite.T(), p.Title, suite.dummy.Title, "unexpected Title")
		assert.Equal(suite.T(), p.Content, suite.dummy.Content, "unexpected Content")
		assert.Equal(suite.T(), p.UserID, suite.dummy.UserID, "unexpected UserID")
		// TODO: research about how to mock parent table
		assert.Equal(suite.T(), p.Username, "", "unexpected Username")
		assert.Equal(suite.T(), p.Publish, suite.dummy.Publish, "unexpected Publish")
		assert.Equal(suite.T(), p.EyeCatchImgName, suite.dummy.EyeCatchImgName, "unexpected EyeCatchImgName")
	})
}

func (suite *ProcedureDaoTestSuite) TestIndex() {
	suite.Run("get procedures with no queries", func() {
		suite.mock.ExpectQuery(
			regexp.QuoteMeta(
				`SELECT procedures.*, users.username FROM "procedures" inner join users ON users.id=procedures.user_id WHERE (title LIKE $1 AND user_id = $2) AND "procedures"."deleted_at" IS NULL ORDER BY updated_at DESC LIMIT 10`),
		).WillReturnRows(suite.mock.NewRows([]string{"id", "title", "content", "user_id", "users.username", "publish", "eye_catch_img_name"}).
			AddRow(suite.dummy.ID, suite.dummy.Title, suite.dummy.Content, suite.dummy.UserID, faker.Username(), suite.dummy.Publish, suite.dummy.EyeCatchImgName))

		ps, err := Index(suite.TestDB, 10, 0, "", suite.dummy.UserID)
		require.NoError(suite.T(), err)
		assert.NotNil(suite.T(), ps)
	})
}

func (suite *ProcedureDaoTestSuite) TestPublicIndex() {
	suite.Run("get procedures with no queries", func() {
		suite.mock.ExpectQuery(
			regexp.QuoteMeta(
				`SELECT procedures.*, users.username FROM "procedures" inner join users ON users.id=procedures.user_id WHERE (title LIKE $1 AND publish = $2) AND "procedures"."deleted_at" IS NULL ORDER BY updated_at DESC LIMIT 10`),
		).WillReturnRows(suite.mock.NewRows([]string{"id", "title", "content", "user_id", "users.username", "publish", "eye_catch_img_name"}).
			AddRow(suite.dummy.ID, suite.dummy.Title, suite.dummy.Content, suite.dummy.UserID, faker.Username(), suite.dummy.Publish, suite.dummy.EyeCatchImgName))

		ps, err := PublicIndex(suite.TestDB, 10, 0, "")
		require.NoError(suite.T(), err)
		assert.NotNil(suite.T(), ps)
	})
}

func (suite *ProcedureDaoTestSuite) TestCreate() {
	suite.Run("create a procedure", func() {
		rows := sqlmock.NewRows([]string{"id"}).AddRow(suite.dummy.ID)
		suite.mock.ExpectBegin()
		suite.mock.ExpectQuery(
			regexp.QuoteMeta(
				`INSERT INTO "procedures" ("created_at","updated_at","deleted_at","title","content","user_id","publish","eye_catch_img_name") VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING "id"`),
		).WillReturnRows(rows)
		suite.mock.ExpectCommit()

		procedure := &Procedure{
			Title:           suite.dummy.Title,
			Content:         suite.dummy.Content,
			UserID:          suite.dummy.UserID,
			Publish:         suite.dummy.Publish,
			EyeCatchImgName: suite.dummy.EyeCatchImgName,
		}
		err := procedure.Save(suite.TestDB)

		require.NoError(suite.T(), err)
		assert.Equal(suite.T(), procedure.ID, suite.dummy.ID, "unexpected ID")
		assert.Equal(suite.T(), procedure.Title, suite.dummy.Title, "unexpected Title")
		assert.Equal(suite.T(), procedure.Content, suite.dummy.Content, "unexpected Content")
		assert.Equal(suite.T(), procedure.UserID, suite.dummy.UserID, "unexpected UserID")
		assert.Equal(suite.T(), procedure.Publish, suite.dummy.Publish, "unexpected Publish")
		assert.Equal(suite.T(), procedure.EyeCatchImgName, suite.dummy.EyeCatchImgName, "unexpected EyeCatchImgName")
	})
}
func (suite *ProcedureDaoTestSuite) TestUpdate() {
	suite.Run("update a procedure", func() {
		suite.mock.ExpectBegin()
		suite.mock.ExpectExec(
			regexp.QuoteMeta(`UPDATE "procedures" SET "created_at"=$1,"updated_at"=$2,"deleted_at"=$3,"title"=$4,"content"=$5,"user_id"=$6,"publish"=$7,"eye_catch_img_name"=$8 WHERE "id" = $9 AND "procedures"."deleted_at" IS NULL`),
		).WillReturnResult(sqlmock.NewResult(int64(suite.dummy.ID), 1))
		suite.mock.ExpectCommit()

		procedure := &Procedure{
			ID:              suite.dummy.ID,
			Title:           faker.Word(),
			Content:         faker.Sentence(),
			Publish:         suite.dummy.Publish,
			EyeCatchImgName: suite.dummy.EyeCatchImgName,
		}
		err := procedure.Update(suite.TestDB)

		require.NoError(suite.T(), err)
		assert.Equal(suite.T(), procedure.ID, suite.dummy.ID, "unexpected ID")
		assert.NotEqual(suite.T(), procedure.Title, suite.dummy.Title, "unexpected Title")
		assert.NotEqual(suite.T(), procedure.Content, suite.dummy.Content, "unexpected Content")
		assert.Empty(suite.T(), procedure.UserID, "unexpected UserID")
		assert.Equal(suite.T(), procedure.Publish, suite.dummy.Publish, "unexpected Publish")
		assert.Equal(suite.T(), procedure.EyeCatchImgName, suite.dummy.EyeCatchImgName, "unexpected EyeCatchImgName")
	})
}

func (suite *ProcedureDaoTestSuite) TestPartialUpdate() {
	suite.Run("partial update 'content' column of a procedure", func() {
		suite.mock.ExpectBegin()
		suite.mock.ExpectExec(
			regexp.QuoteMeta(`UPDATE "procedures" SET "updated_at"=$1,"content"=$2,"publish"=$3 WHERE id IN ($4) AND "id" = $5 AND "procedures"."deleted_at" IS NULL`),
		).WillReturnResult(sqlmock.NewResult(int64(suite.dummy.ID), 1))
		suite.mock.ExpectCommit()

		procedure := &Procedure{
			ID:      suite.dummy.ID,
			Content: faker.Sentence(),
			Publish: suite.dummy.Publish,
		}
		err := procedure.PartialUpdate(suite.TestDB)

		require.NoError(suite.T(), err)
		assert.Equal(suite.T(), procedure.ID, suite.dummy.ID, "unexpected ID")
		assert.Empty(suite.T(), procedure.Title, "unexpected Title")
		assert.NotEqual(suite.T(), procedure.Content, suite.dummy.Content, "unexpected Content")
		assert.Empty(suite.T(), procedure.UserID, "unexpected UserID")
		assert.Equal(suite.T(), procedure.Publish, suite.dummy.Publish, "unexpected Publish")
	})
}
func (suite *ProcedureDaoTestSuite) TestDelete() {
	suite.Run("delete a procedure", func() {
		suite.mock.ExpectBegin()
		suite.mock.ExpectExec(
			`UPDATE "procedures" SET "deleted_at"=(.*)WHERE "procedures"."id" = (.*)"procedures"."deleted_at" IS NULL`).WillReturnResult(sqlmock.NewResult(int64(suite.dummy.ID), 1))
		suite.mock.ExpectCommit()

		procedure := &Procedure{
			ID:              suite.dummy.ID,
			Title:           suite.dummy.Title,
			Content:         suite.dummy.Content,
			UserID:          suite.dummy.UserID,
			Publish:         suite.dummy.Publish,
			EyeCatchImgName: suite.dummy.EyeCatchImgName,
		}
		err := procedure.Delete(suite.TestDB)

		require.NoError(suite.T(), err)
		assert.Equal(suite.T(), procedure.ID, suite.dummy.ID, "unexpected ID")
		assert.Equal(suite.T(), procedure.Title, suite.dummy.Title, "unexpected Title")
		assert.Equal(suite.T(), procedure.Content, suite.dummy.Content, "unexpected Content")
		assert.Equal(suite.T(), procedure.UserID, suite.dummy.UserID, "unexpected UserID")
		assert.Equal(suite.T(), procedure.Publish, suite.dummy.Publish, "unexpected Publish")
		assert.Equal(suite.T(), procedure.EyeCatchImgName, suite.dummy.EyeCatchImgName, "unexpected EyeCatchImgName")
	})
}
