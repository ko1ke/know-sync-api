package steps

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

type StepDaoTestSuite struct {
	suite.Suite
	TestDB *gorm.DB
	mock   sqlmock.Sqlmock
	dummy  Step
}

// set up test to gorm send queries to mock
func (suite *StepDaoTestSuite) SetupTest() {
	db, mock, _ := sqlmock.New()
	suite.mock = mock
	suite.TestDB, _ = gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
}

func (suite *StepDaoTestSuite) TearDownTest() {
	db, _ := suite.TestDB.DB()
	db.Close()
}

func (suite *StepDaoTestSuite) AfterTest(_, _ string) {
	require.NoError(suite.T(), suite.mock.ExpectationsWereMet())
}

// Run test suite
func TestStepDaoTestSuite(t *testing.T) {
	suite.Run(t, new(StepDaoTestSuite))
}

func (suite *StepDaoTestSuite) BeforeTest(suiteName string, testName string) {

	suite.dummy = Step{
		ID:          rand_utils.MakeRandomUInt(100),
		Content:     faker.Sentence(),
		ImgName:     faker.Word() + ".jpeg",
		ProcedureID: rand_utils.MakeRandomUInt(100),
	}
	// spew.Dump(suite.dummy)
}

func (suite *StepDaoTestSuite) TestGet() {
	suite.Run("get a Step", func() {
		suite.mock.ExpectQuery(
			regexp.QuoteMeta(`SELECT * FROM "steps" WHERE id = $1 AND "steps"."deleted_at" IS NULL AND "steps"."id" = $2 ORDER BY "steps"."id" LIMIT 1`),
		).WithArgs(suite.dummy.ID, suite.dummy.ID).
			WillReturnRows(suite.mock.NewRows([]string{"id", "content", "procedure_id", "img_name"}).
				AddRow(suite.dummy.ID, suite.dummy.Content, suite.dummy.ProcedureID, suite.dummy.ImgName))

		step := &Step{
			ID: suite.dummy.ID,
		}
		err := step.Get(suite.TestDB)
		require.NoError(suite.T(), err)
		assert.Equal(suite.T(), step.ID, suite.dummy.ID, "unexpected ID")
		assert.Equal(suite.T(), step.Content, suite.dummy.Content, "unexpected Content")
		assert.Equal(suite.T(), step.ImgName, suite.dummy.ImgName, "unexpected ImgName")
		assert.Equal(suite.T(), step.ProcedureID, suite.dummy.ProcedureID, "unexpected ProcedureID")
	})
}

func (suite *StepDaoTestSuite) TestIndex() {
	suite.Run("get steps", func() {
		suite.mock.ExpectQuery(
			regexp.QuoteMeta(
				`SELECT * FROM "steps" WHERE procedure_id = $1 AND "steps"."deleted_at" IS NULL`),
		).WillReturnRows(suite.mock.NewRows([]string{"id", "content", "procedure_id", "img_name"}).
			AddRow(suite.dummy.ID, suite.dummy.Content, suite.dummy.ProcedureID, suite.dummy.ImgName))

		ps, err := Index(suite.TestDB, suite.dummy.ProcedureID)
		require.NoError(suite.T(), err)
		assert.NotNil(suite.T(), ps)
	})
}

func (suite *StepDaoTestSuite) TestCreate() {
	suite.Run("create a step", func() {
		rows := sqlmock.NewRows([]string{"id"}).AddRow(suite.dummy.ID)
		suite.mock.ExpectBegin()
		suite.mock.ExpectQuery(
			regexp.QuoteMeta(
				`INSERT INTO "steps" ("created_at","updated_at","deleted_at","content","img_name","procedure_id") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "id"`),
		).WillReturnRows(rows)
		suite.mock.ExpectCommit()

		step := &Step{
			Content:     suite.dummy.Content,
			ImgName:     suite.dummy.ImgName,
			ProcedureID: suite.dummy.ProcedureID,
		}
		err := step.Save(suite.TestDB)

		require.NoError(suite.T(), err)
		assert.Equal(suite.T(), step.ID, suite.dummy.ID, "unexpected ID")
		assert.Equal(suite.T(), step.Content, suite.dummy.Content, "unexpected Content")
		assert.Equal(suite.T(), step.ImgName, suite.dummy.ImgName, "unexpected ImgName")
		assert.Equal(suite.T(), step.ProcedureID, suite.dummy.ProcedureID, "unexpected ProcedureID")
	})
}

func (suite *StepDaoTestSuite) TestBulkCreate() {
	suite.Run("create steps", func() {
		suite.mock.ExpectBegin()
		rows := sqlmock.NewRows([]string{"id"}).AddRow(suite.dummy.ID).AddRow(rand_utils.MakeRandomUInt(100))
		suite.mock.ExpectQuery(
			regexp.QuoteMeta(
				`INSERT INTO "steps" ("created_at","updated_at","deleted_at","content","img_name","procedure_id") VALUES ($1,$2,$3,$4,$5,$6),($7,$8,$9,$10,$11,$12) ON CONFLICT ("id") DO UPDATE SET "updated_at"=$13,"deleted_at"="excluded"."deleted_at","content"="excluded"."content","img_name"="excluded"."img_name","procedure_id"="excluded"."procedure_id" RETURNING "id"`),
		).WillReturnRows(rows)
		suite.mock.ExpectCommit()

		step1 := Step{
			Content:     suite.dummy.Content,
			ImgName:     suite.dummy.ImgName,
			ProcedureID: suite.dummy.ProcedureID,
		}
		step2 := Step{
			Content:     faker.Sentence(),
			ImgName:     faker.Word() + ".jpeg",
			ProcedureID: suite.dummy.ProcedureID,
		}

		var ss []Step
		ss = append(ss, step1)
		ss = append(ss, step2)
		ss, err := BulkCreate(suite.TestDB, ss)

		require.NoError(suite.T(), err)
		assert.NotNil(suite.T(), ss)
	})
}
func (suite *StepDaoTestSuite) TestUpdate() {
	suite.Run("update a step", func() {
		suite.mock.ExpectBegin()
		suite.mock.ExpectExec(
			regexp.QuoteMeta(`UPDATE "steps" SET "created_at"=$1,"updated_at"=$2,"deleted_at"=$3,"content"=$4,"img_name"=$5,"procedure_id"=$6 WHERE "id" = $7 AND "steps"."deleted_at" IS NULL`),
		).WillReturnResult(sqlmock.NewResult(int64(suite.dummy.ID), 1))
		suite.mock.ExpectCommit()

		step := &Step{
			ID:      suite.dummy.ID,
			Content: faker.Sentence(),
		}
		err := step.Update(suite.TestDB)

		require.NoError(suite.T(), err)
		assert.Equal(suite.T(), step.ID, suite.dummy.ID, "unexpected ID")
		assert.NotEqual(suite.T(), step.Content, suite.dummy.Content, "unexpected Content")
		assert.Empty(suite.T(), step.ImgName, "unexpected ImgName")
		assert.Empty(suite.T(), step.ProcedureID, "unexpected ProcedureID")
	})
}

func (suite *StepDaoTestSuite) TestDelete() {
	suite.Run("delete a step", func() {
		suite.mock.ExpectBegin()
		suite.mock.ExpectExec(
			`UPDATE "steps" SET "deleted_at"=(.*)WHERE "steps"."id" = (.*)"steps"."deleted_at" IS NULL`).WillReturnResult(sqlmock.NewResult(int64(suite.dummy.ID), 1))
		suite.mock.ExpectCommit()

		step := &Step{
			ID:          suite.dummy.ID,
			Content:     suite.dummy.Content,
			ImgName:     suite.dummy.ImgName,
			ProcedureID: suite.dummy.ProcedureID,
		}
		err := step.Delete(suite.TestDB)

		require.NoError(suite.T(), err)
		assert.Equal(suite.T(), step.ID, suite.dummy.ID, "unexpected ID")
		assert.Equal(suite.T(), step.Content, suite.dummy.Content, "unexpected Content")
		assert.Equal(suite.T(), step.ImgName, suite.dummy.ImgName, "unexpected ImgName")
		assert.Equal(suite.T(), step.ProcedureID, suite.dummy.ProcedureID, "unexpected ProcedureID")
	})
}
func (suite *StepDaoTestSuite) TestBulkDeleteByProcedureID() {
	suite.Run("delete steps by procedureID", func() {
		procedureId := uint(1)
		sqlmock.NewRows([]string{"id", "procedure_id"}).AddRow(1, procedureId).AddRow(2, procedureId)
		suite.mock.ExpectBegin()
		suite.mock.ExpectExec(
			regexp.QuoteMeta(
				`UPDATE "steps" SET "deleted_at"=$1 WHERE procedure_id = $2 AND "steps"."deleted_at" IS NULL`)).WillReturnResult(sqlmock.NewResult(int64(2), 2))
		suite.mock.ExpectCommit()
		err := BulkDeleteByProcedureId(suite.TestDB, procedureId)
		require.NoError(suite.T(), err)
	})
}
