package questionAndAnswerServices

import (
	"codeSandbox/model/vo"
	"codeSandbox/utils/global"
	"encoding/json"
	log "github.com/sirupsen/logrus"

	"os"
)

const (
	FILE_NAME = "utils/qas/qas.json"
)

type QuestionAndAnswerService struct {
}

var QuestionAndAnswerServiceInstance QuestionAndAnswerService

func (service *QuestionAndAnswerService) GetAllQuestionAndAnswers() (int, []vo.QuestionAndAnswer) {
	data, err := os.ReadFile(FILE_NAME)
	if err != nil {
		log.Errorf("Read %v fail:%v", FILE_NAME, err)
		return global.SYSTEM_ERROR, nil
	}
	questionAndAnswers := make([]vo.QuestionAndAnswer, 0, 0)
	err = json.Unmarshal(data, &questionAndAnswers)
	if err != nil {
		log.Errorf("json.Unmarshalfail:%v", err)
		return global.SYSTEM_ERROR, nil
	}
	return global.SUCCESS, questionAndAnswers
}
