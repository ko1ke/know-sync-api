package users

import (
	"regexp"
	"testing"

	"github.com/ko1ke/know-sync-api/cmd/utils/rand_utils"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type UserDaoTestSuite struct {
	suite.Suite
	TestDB *gorm.DB
	mock   sqlmock.Sqlmock
	dummy  User
}

// set up test to gorm send queries to mock
func (suite *UserDaoTestSuite) SetupTest() {
	db, mock, _ := sqlmock.New()
	suite.mock = mock
	suite.TestDB, _ = gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
}

func (suite *UserDaoTestSuite) TearDownTest() {
	db, _ := suite.TestDB.DB()
	db.Close()
}

func (suite *UserDaoTestSuite) AfterTest(_, _ string) {
	require.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

// Run test suite
func TestUserDaoTestSuite(t *testing.T) {
	suite.Run(t, new(UserDaoTestSuite))
}

func (suite *UserDaoTestSuite) BeforeTest(suiteName string, testName string) {

	suite.dummy = User{
		ID:       rand_utils.MakeRandomUInt(100),
		Username: faker.Name(),
		Email:    faker.Email(),
		Password: faker.Password(),
	}
	// spew.Dump(suite.dummy)
}

func (suite *UserDaoTestSuite) TestGet() {
	suite.Run("get a procedure", func() {
		suite.mock.ExpectQuery(
			regexp.QuoteMeta(`SELECT * FROM "users" WHERE id = $1 AND "users"."deleted_at" IS NULL AND "users"."id" = $2 ORDER BY "users"."id" LIMIT 1`),
		).WithArgs(suite.dummy.ID, suite.dummy.ID).
			WillReturnRows(suite.mock.NewRows([]string{"id", "username", "email", "password"}).
				AddRow(suite.dummy.ID, suite.dummy.Username, suite.dummy.Email, suite.dummy.Password))

		procedure := &User{
			ID: suite.dummy.ID,
		}
		err := procedure.Get(suite.TestDB)
		require.NoError(suite.T(), err)
		assert.Equal(suite.T(), procedure.ID, suite.dummy.ID, "unexpected ID")
		assert.Equal(suite.T(), procedure.Username, suite.dummy.Username, "unexpected Username")
		assert.Equal(suite.T(), procedure.Email, suite.dummy.Email, "unexpected Email")
		assert.Equal(suite.T(), procedure.Password, suite.dummy.Password, "unexpected Password")
	})
}

func (suite *UserDaoTestSuite) TestCreate() {
	suite.Run("create a procedure", func() {
		rows := sqlmock.NewRows([]string{"id"}).AddRow(suite.dummy.ID)
		suite.mock.ExpectBegin()
		suite.mock.ExpectQuery(
			regexp.QuoteMeta(
				`INSERT INTO "users" ("created_at",` +
					`"updated_at","deleted_at","username",` +
					`"email","password") VALUES ($1,$2,$3,$4,$5,$6) ` +
					`RETURNING "id"`),
		).WillReturnRows(rows)
		suite.mock.ExpectCommit()

		procedure := &User{
			Username: suite.dummy.Username,
			Email:    suite.dummy.Email,
			Password: suite.dummy.Password,
		}
		err := procedure.Save(suite.TestDB)

		require.NoError(suite.T(), err)
		assert.Equal(suite.T(), procedure.ID, suite.dummy.ID, "unexpected ID")
		assert.Equal(suite.T(), procedure.Username, suite.dummy.Username, "unexpected Username")
		assert.Equal(suite.T(), procedure.Email, suite.dummy.Email, "unexpected Email")
		assert.Equal(suite.T(), procedure.Password, suite.dummy.Password, "unexpected Password")
	})
}
func (suite *UserDaoTestSuite) TestDelete() {
	suite.Run("delete a procedure", func() {
		suite.mock.ExpectBegin()
		suite.mock.ExpectExec(
			`UPDATE "users" SET "deleted_at"=(.*)WHERE "users"."id" = (.*)"users"."deleted_at" IS NULL`).WillReturnResult(sqlmock.NewResult(int64(suite.dummy.ID), 1))
		suite.mock.ExpectCommit()

		procedure := &User{
			ID:       suite.dummy.ID,
			Username: suite.dummy.Username,
			Email:    suite.dummy.Email,
			Password: suite.dummy.Password,
		}
		err := procedure.Delete(suite.TestDB)

		require.NoError(suite.T(), err)
		assert.Equal(suite.T(), procedure.ID, suite.dummy.ID, "unexpected ID")
		assert.Equal(suite.T(), procedure.Username, suite.dummy.Username, "unexpected Username")
		assert.Equal(suite.T(), procedure.Email, suite.dummy.Email, "unexpected Email")
		assert.Equal(suite.T(), procedure.Password, suite.dummy.Password, "unexpected Password")
	})
}
