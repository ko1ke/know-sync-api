package procedures

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ProcedureDaoTestSuite struct {
	suite.Suite
	TestDB *gorm.DB
	mock   sqlmock.Sqlmock
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

// Run test suite
func TestProcedureDaoTestSuite(t *testing.T) {
	suite.Run(t, new(ProcedureDaoTestSuite))
}

func (suite *ProcedureDaoTestSuite) TestCreate() {
	suite.Run("create a procedure", func() {
		newId := uint(1)
		rows := sqlmock.NewRows([]string{"id"}).AddRow(newId)
		suite.mock.ExpectBegin()
		suite.mock.ExpectQuery(
			regexp.QuoteMeta(
				`INSERT INTO "procedures" ("created_at",` +
					`"updated_at","deleted_at","title",` +
					`"content","user_id") VALUES ($1,$2,$3,$4,$5,$6) ` +
					`RETURNING "id"`),
		).WillReturnRows(rows)
		suite.mock.ExpectCommit()
		title := "ルアーのウレタンコーティング"
		content := "ウレタンにドブ漬けする"
		userId := uint(1)

		procedure := &Procedure{
			Title:   title,
			Content: content,
			UserID:  userId,
		}
		err := procedure.Save(suite.TestDB)

		if err != nil {
			suite.Fail("Error:" + err.Error())
		}
		assert.Equal(suite.T(), procedure.ID, newId, "unexpected ID")
		assert.Equal(suite.T(), procedure.Title, title, "unexpected Title")
		assert.Equal(suite.T(), procedure.Content, content, "unexpected Content")
		assert.Equal(suite.T(), procedure.UserID, userId, "unexpected UserID")
	})
}
