package mailer

import (
	"fmt"

	"gopkg.in/gomail.v2"

	"zoomment-server/internal/config"
	"zoomment-server/internal/logger"
)

// Mailer handles sending emails
type Mailer struct {
	dialer    *gomail.Dialer
	from      string
	brandName string
	dashboardURL string
}

// New creates a new Mailer instance
func New(cfg *config.Config) *Mailer {
	dialer := gomail.NewDialer(
		cfg.BotEmail.Host,
		cfg.BotEmail.Port,
		cfg.BotEmail.Address,
		cfg.BotEmail.Password,
	)

	return &Mailer{
		dialer:    dialer,
		from:      cfg.BotEmail.Address,
		brandName: cfg.BrandName,
		dashboardURL: cfg.DashboardURL,
	}
}

// SendMagicLink sends a magic link email for authentication
func (m *Mailer) SendMagicLink(email, token string) error {
	if m.from == "" {
		logger.Warn("Email not configured, skipping magic link email")
		return nil
	}

	link := fmt.Sprintf("%s/dashboard?zoommentToken=%s", m.dashboardURL, token)

	html := generateTemplate(TemplateData{
		BrandName:    m.brandName,
		DashboardURL: m.dashboardURL,
		Introduction: fmt.Sprintf("Click the link below to sign in to your %s dashboard.", m.brandName),
		ButtonText:   fmt.Sprintf("Sign in to %s", m.brandName),
		ButtonURL:    link,
		Epilogue:     "If you did not make this request, you can safely ignore this email.",
	})

	msg := gomail.NewMessage()
	msg.SetHeader("From", fmt.Sprintf("%s <%s>", m.brandName, m.from))
	msg.SetHeader("To", email)
	msg.SetHeader("Subject", fmt.Sprintf("Sign in to %s", m.brandName))
	msg.SetBody("text/html", html)

	if err := m.dialer.DialAndSend(msg); err != nil {
		logger.Error(err, "Failed to send magic link email")
		return err
	}

	logger.Info("Magic link email sent to " + email)
	return nil
}

// SendEmailVerification sends email verification link for guest comments
func (m *Mailer) SendEmailVerification(email, token, pageURL string) error {
	if m.from == "" {
		logger.Warn("Email not configured, skipping verification email")
		return nil
	}

	link := fmt.Sprintf("%s?zoommentToken=%s", pageURL, token)

	html := generateTemplate(TemplateData{
		BrandName:    m.brandName,
		DashboardURL: m.dashboardURL,
		Introduction: "Please confirm your email address to be able to manage your comment.",
		ButtonText:   "Confirm",
		ButtonURL:    link,
		Epilogue:     "If you did not make this request, you can safely ignore this email.",
	})

	msg := gomail.NewMessage()
	msg.SetHeader("From", fmt.Sprintf("%s <%s>", m.brandName, m.from))
	msg.SetHeader("To", email)
	msg.SetHeader("Subject", "You have added a comment!")
	msg.SetBody("text/html", html)

	if err := m.dialer.DialAndSend(msg); err != nil {
		logger.Error(err, "Failed to send verification email")
		return err
	}

	logger.Info("Verification email sent to " + email)
	return nil
}

// SendCommentNotification notifies site owner about a new comment
func (m *Mailer) SendCommentNotification(ownerEmail string, comment CommentData) error {
	if m.from == "" {
		logger.Warn("Email not configured, skipping notification email")
		return nil
	}

	html := generateTemplate(TemplateData{
		BrandName:    m.brandName,
		DashboardURL: m.dashboardURL,
		Introduction: fmt.Sprintf(`
			<p>You have a new comment!</p>
			<div style="font-size: 14px; line-height: 27px; margin-top: 10px;"> 
				<div><b>User:</b> %s</div>
				<div><b>Date:</b> %s</div>
				<div><b>Page:</b> %s</div>
				<div><b>Comment:</b> %s</div>
			</div>
		`, comment.Author, comment.Date, comment.PageURL, comment.Body),
		ButtonText: "Sign in to manage comments",
		ButtonURL:  m.dashboardURL + "/auth",
		Epilogue:   "",
	})

	msg := gomail.NewMessage()
	msg.SetHeader("From", fmt.Sprintf("%s <%s>", m.brandName, m.from))
	msg.SetHeader("To", ownerEmail)
	msg.SetHeader("Subject", "You have a new comment!")
	msg.SetBody("text/html", html)

	if err := m.dialer.DialAndSend(msg); err != nil {
		logger.Error(err, "Failed to send notification email")
		return err
	}

	logger.Info("Notification email sent to " + ownerEmail)
	return nil
}

// CommentData holds data for comment notification email
type CommentData struct {
	Author  string
	Date    string
	PageURL string
	Body    string
}

