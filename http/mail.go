package http

import (
	"net/http"
	"strings"

	"github.com/mail-provider/config"
	"github.com/peng19940915/smtp"
	"github.com/toolkits/web/param"
)

func configProcRoutes() {

	http.HandleFunc("/sender/mail", func(w http.ResponseWriter, r *http.Request) {
		cfg := config.Config()
		token := param.String(r, "token", "")
		if cfg.Http.Token != token {
			http.Error(w, "no privilege", http.StatusForbidden)
			return
		}

		tos := param.MustString(r, "tos")
		subject := param.MustString(r, "subject")
		content := param.MustString(r, "content")
		tos = strings.Replace(tos, ",", ";", -1)
		// 强制关闭匿名邮件
		s := smtp.NewSMTP(cfg.Smtp.Addr, cfg.Smtp.Username, cfg.Smtp.Password, cfg.Smtp.TLS, false, cfg.Smtp.SkipVerify)
		err := s.SendMail(cfg.Smtp.From, tos, subject, content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			http.Error(w, "success", http.StatusOK)
		}
	})

}
