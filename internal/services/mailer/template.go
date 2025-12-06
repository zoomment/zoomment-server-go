package mailer

import "fmt"

// TemplateData holds data for email templates
type TemplateData struct {
	BrandName    string
	DashboardURL string
	Introduction string
	ButtonText   string
	ButtonURL    string
	Epilogue     string
}

// generateTemplate creates the HTML email template
// This is the same template as your Node.js version
func generateTemplate(data TemplateData) string {
	buttonHTML := ""
	if data.ButtonURL != "" {
		buttonHTML = fmt.Sprintf(`
			<br/>
			<table role="presentation" border="0" cellpadding="0" cellspacing="0" class="btn btn-primary">
				<tbody>
					<tr>
						<td align="center">
							<table role="presentation" border="0" cellpadding="0" cellspacing="0">
								<tbody>
									<tr>
										<td>
											<a href="%s" target="_blank">%s</a> 
										</td>
									</tr>
								</tbody>
							</table>
						</td>
					</tr>
				</tbody>
			</table>
			<br/>
		`, data.ButtonURL, data.ButtonText)
	}

	return fmt.Sprintf(`
<!doctype html>
<html lang="en">
<head>
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
	<title>%s</title>
	<style media="all" type="text/css">
		body {
			font-family: Helvetica, sans-serif;
			-webkit-font-smoothing: antialiased;
			font-size: 14px;
			line-height: 1.3;
			-ms-text-size-adjust: 100%%;
			-webkit-text-size-adjust: 100%%;
			background-color: #f5f5f5;
			margin: 0;
			padding: 0;
		}
		.container {
			margin: 0 auto !important;
			max-width: 400px;
			padding: 24px 0;
		}
		.main {
			background: #ffffff;
			border: 1px solid #eaebed;
			border-radius: 14px;
			padding: 24px 35px;
		}
		h2 {
			font-size: 18px;
			font-weight: 700;
			margin: 12px 0 16px;
			text-align: center;
		}
		p {
			margin: 0 0 16px;
		}
		.btn-primary a {
			background-color: #1677ff;
			border: solid 2px #1677ff;
			border-radius: 8px;
			color: #ffffff;
			display: inline-block;
			font-size: 16px;
			font-weight: bold;
			padding: 10px 20px;
			text-decoration: none;
		}
		.btn-primary a:hover {
			background-color: #4096ff;
		}
		.footer {
			text-align: center;
			padding: 24px;
			color: #9a9ea6;
			font-size: 14px;
		}
		.logo {
			text-align: center;
			margin-bottom: 20px;
		}
	</style>
</head>
<body>
	<div class="container">
		<div class="main">
			<div class="logo">
				<img width="50" height="auto" src="%s/email-logo.png" alt="%s" />
				<h2>%s</h2>
			</div>
			<p>%s</p>
			%s
			<p>%s</p>
		</div>
		<div class="footer">
			%s
		</div>
	</div>
</body>
</html>
`, data.BrandName, data.DashboardURL, data.BrandName, data.BrandName, 
   data.Introduction, buttonHTML, data.Epilogue, data.BrandName)
}

