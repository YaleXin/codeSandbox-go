package test

import (
	"codeSandbox/db"
	"testing"
)

func TestDaoPage(t *testing.T) {
	executionDao := db.ExecutionDao{}
	executions, total, err := executionDao.PageExecution(5, 1)
	if err != nil {
		t.Fatalf("err:%v", err)
	}
	executions2, total2, err := executionDao.PageExecutionByUserId(21, 5, 1)
	if err != nil {
		t.Fatalf("err:%v", err)
	}
	t.Logf("total:%v", total)
	t.Logf("executions[0].User:%v", executions[0].User)

	t.Logf("total2:%v", total2)
	t.Logf("executions2[0].User:%v", executions2[0].User)
}
