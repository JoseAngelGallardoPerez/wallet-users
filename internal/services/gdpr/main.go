package gdpr

import (
	"fmt"
	"html"
	"time"

	"github.com/Confialink/wallet-pkg-utils/timefmt"

	"github.com/Confialink/wallet-users/internal/services/customization"
	"github.com/Confialink/wallet-users/internal/services/pdf"
	"github.com/Confialink/wallet-users/internal/services/syssettings"
)

type Service struct {
	sysSettingsService *syssettings.SysSettings
}

func NewService(sysSettingsService *syssettings.SysSettings) *Service {
	return &Service{sysSettingsService: sysSettingsService}
}

// GdprHtmlBytes requests a GDPR policy from the Customization service and generates a byte slice for a PDF file
func (s *Service) GdprHtmlBytes() ([]byte, error) {
	gdpr, err := customization.GetCustomizationByKey("gdpr")
	if err != nil {
		return nil, err
	}

	settings, err := s.sysSettingsService.GetTimeSettings()
	if err != nil {
		return nil, err
	}

	l, err := time.LoadLocation(settings.Timezone)
	if err != nil {
		return nil, err
	}
	_, o := time.Now().In(l).Zone()
	h := o / 3600
	var sep string
	if h > 0 {
		sep = "+"
	}
	z := fmt.Sprintf("UTC%s%d", sep, h)
	t := fmt.Sprintf("%s (%s)", timefmt.Format(time.Now(), settings.DateTimeFormat, settings.Timezone), z)

	return pdf.HtmlToPdfBytes(html.UnescapeString(`<!doctype html><html><body>` +
		`<h1>` + gdpr.Label + `</h1><br>` + gdpr.Value + `<br><h4>` + t + `</h4>` + `</body></html>`))
}
