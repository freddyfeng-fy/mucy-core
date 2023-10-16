package markdown

import (
	"github.com/88250/lute"
	"github.com/freddyfeng-fy/mucy-core/html"
	"github.com/freddyfeng-fy/mucy-core/utils/strs"
	"sync"
)

var (
	engine *lute.Lute
	once   sync.Once
)

func getEngine() *lute.Lute {
	once.Do(func() {
		engine = lute.New(func(lute *lute.Lute) {
			lute.SetToC(true)
			lute.SetSanitize(true)
			lute.SetGFMTaskListItem(true)
		})
	})
	return engine
}

func ToHTML(markdownStr string) string {
	if strs.IsBlank(markdownStr) {
		return ""
	}
	return getEngine().MarkdownStr("", markdownStr)
}

func GetSummary(markdownStr string, summaryLen int) string {
	htmlStr := ToHTML(markdownStr)
	return html.GetSummary(htmlStr, summaryLen)
}
