package weixin

import (
	"github.com/kkserver/kk-lib/kk/app"
)

type WXTicketTaskResult struct {
	app.Result
	Ticket *Ticket `json:"ticket,omitempty"`
}

type WXTicketTask struct {
	app.Task
	Result WXTicketTaskResult
}

func (task *WXTicketTask) GetResult() interface{} {
	return &task.Result
}

func (task *WXTicketTask) GetInhertType() string {
	return "weixin"
}

func (task *WXTicketTask) GetClientName() string {
	return "Ticket"
}
