package procedures

import (
	"know-sync-api/utils/rand_utils"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ProcedureDaoTestSuite struct {
	suite.Suite
	TestDB *gorm.DB
	mock   sqlmock.Sqlmock
	dummy  Procedure
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
		ID:      rand_utils.MakeRandomUInt(100),
		Title:   faker.Word(),
		Content: faker.Sentence(),
		UserID:  rand_utils.MakeRandomUInt(100),
	}
	// spew.Dump(suite.dummy)
}

func (suite *ProcedureDaoTestSuite) TestIndex() {
	suite.Run("get procedures with no queries", func() {
		suite.mock.ExpectQuery(
			regexp.QuoteMeta(
				`SELECT * FROM "procedures" WHERE "procedures"."deleted_at" IS NULL LIMIT 10`),
		).WillReturnRows(suite.mock.NewRows([]string{"id", "title", "content", "userID"}).
			AddRow(suite.dummy.ID, suite.dummy.Title, suite.dummy.Content, suite.dummy.UserID))

		ps, err := Index(suite.TestDB, 10, 0)
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
				`INSERT INTO "procedures" ("created_at",` +
					`"updated_at","deleted_at","title",` +
					`"content","user_id") VALUES ($1,$2,$3,$4,$5,$6) ` +
					`RETURNING "id"`),
		).WillReturnRows(rows)
		suite.mock.ExpectCommit()

		procedure := &Procedure{
			Title:   suite.dummy.Title,
			Content: suite.dummy.Content,
			UserID:  suite.dummy.UserID,
		}
		err := procedure.Save(suite.TestDB)

		require.NoError(suite.T(), err)
		assert.Equal(suite.T(), procedure.ID, suite.dummy.ID, "unexpected ID")
		assert.Equal(suite.T(), procedure.Title, suite.dummy.Title, "unexpected Title")
		assert.Equal(suite.T(), procedure.Content, suite.dummy.Content, "unexpected Content")
		assert.Equal(suite.T(), procedure.UserID, suite.dummy.UserID, "unexpected UserID")
	})
}

func (suite *ProcedureDaoTestSuite) TestDelete() {
	suite.Run("delete a procedure", func() {
		suite.mock.ExpectBegin()
		suite.mock.ExpectExec(
			`UPDATE "procedures" SET "deleted_at"=(.*)WHERE "procedures"."id" = (.*)"procedures"."deleted_at" IS NULL`).WillReturnResult(sqlmock.NewResult(int64(suite.dummy.ID), 1))
		suite.mock.ExpectCommit()

		procedure := &Procedure{
			ID:      suite.dummy.ID,
			Title:   suite.dummy.Title,
			Content: suite.dummy.Content,
			UserID:  suite.dummy.UserID,
		}
		err := procedure.Delete(suite.TestDB)

		require.NoError(suite.T(), err)
		assert.Equal(suite.T(), procedure.ID, suite.dummy.ID, "unexpected ID")
		assert.Equal(suite.T(), procedure.Title, suite.dummy.Title, "unexpected Title")
		assert.Equal(suite.T(), procedure.Content, suite.dummy.Content, "unexpected Content")
		assert.Equal(suite.T(), procedure.UserID, suite.dummy.UserID, "unexpected UserID")
	})
}
