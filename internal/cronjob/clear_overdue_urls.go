package cronjob

import (
	"github.com/vfilipovsky/url-shortener/internal/service"
	"github.com/vfilipovsky/url-shortener/pkg/logger"
)

type clearOverdueUrls struct {
	urlService service.Url
}

func NewClearOverDueUrls(urlService service.Url) *clearOverdueUrls {
	return &clearOverdueUrls{urlService: urlService}
}

func (job *clearOverdueUrls) GetName() string {
	return "Clear overdue urls"
}

func (job *clearOverdueUrls) Run() {
	logger.Infof("[%s] Run", job.GetName())

	err := job.urlService.RemoveOverdue()

	if err != nil {
		logger.Errorf("[%s] Error: %s", job.GetName(), err.Error())
		return
	}

	logger.Infof("[%s] Finish", job.GetName())
}
