package handler

import (
	"fmt"
	"io"
	"text/template"

	"github.com/jumpserver/koko/pkg/i18n"
	"github.com/jumpserver/koko/pkg/jms-sdk-go/model"
	"github.com/jumpserver/koko/pkg/logger"
	"github.com/jumpserver/koko/pkg/utils"
)

type MenuItem struct {
	instruct string
	helpText string
}

type Menu []MenuItem

type ColorMeta struct {
	GreenBoldColor string
	ColorEnd       string
}

func (h *InteractiveHandler) displayBanner(sess io.ReadWriter, user string, termConf *model.TerminalConfig) {
	lang := i18n.NewLang(h.i18nLang)
	defaultTitle := utils.WrapperTitle(lang.T("Welcome to use JumpServer open source fortress system"))
	menu := Menu{
		{instruct: lang.T("part IP, Hostname, Comment"), helpText: lang.T("to search login if unique")},
		{instruct: lang.T("/ + IP, Hostname, Comment"), helpText: lang.T("to search, such as: /192.168")},
		{instruct: "p", helpText: lang.T("display the assets you have permission")},
		{instruct: "g", helpText: lang.T("display the node that you have permission")},
		{instruct: "h", helpText: lang.T("display the hosts that you have permission")},
		{instruct: "d", helpText: lang.T("display the databases that you have permission")},
		{instruct: "k", helpText: lang.T("display the kubernetes that you have permission")},
		{instruct: "r", helpText: lang.T("refresh your assets and nodes")},
		{instruct: "s", helpText: lang.T("Chinese-English-Japanese switch")},
		{instruct: "?", helpText: lang.T("print help")},
		{instruct: "q", helpText: lang.T("exit")},
	}

	title := defaultTitle
	if termConf.HeaderTitle != "" {
		title = termConf.HeaderTitle
	}

	prefix := utils.CharClear + utils.CharTab + utils.CharTab
	suffix := utils.CharNewLine + utils.CharNewLine
	welcomeMsg := prefix + utils.WrapperTitle(user+",") + "  " + title + suffix
	_, err := io.WriteString(sess, welcomeMsg)
	if err != nil {
		logger.Errorf("Send to client error, %s", err)
		return
	}
	cm := ColorMeta{GreenBoldColor: "\033[1;32m", ColorEnd: "\033[0m"}
	for i, v := range menu {
		line := fmt.Sprintf(lang.T("\t%d) Enter {{.GreenBoldColor}}%s{{.ColorEnd}} to %s.%s"),
			i+1, v.instruct, v.helpText, "\r\n")
		tmpl := template.Must(template.New("item").Parse(line))
		if err := tmpl.Execute(sess, cm); err != nil {
			logger.Error(err)
		}
	}
}
