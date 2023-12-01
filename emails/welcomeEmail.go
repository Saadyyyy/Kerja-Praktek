package emails

import (
	"os"
	"strconv"

	"github.com/go-gomail/gomail"
)

func SendWelcomeEmail(userEmail, name, verificationToken string) error {
	smtpServer := os.Getenv("SMTP_SERVER")
	smtpPortStr := os.Getenv("SMTP_PORT")
	smtpUsername := os.Getenv("SMTP_USERNAME")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	sender := smtpUsername
	recipient := userEmail
	subject := "Laode Saady Website"
	verificationLink := "http://localhost:8080/verify?token=" + verificationToken
	emailBody := `
    <html>
    <head>
        <link href="https://maxcdn.bootstrapcdn.com/bootstrap/4.5.2/css/bootstrap.min.css" rel="stylesheet">
        <style>
            body {
    font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif;
    background-color: #f5f5f5;
    color: #333;
}

            .container {
    max-width: 600px;
    margin: 0 auto;
    padding: 20px;
    background-color: #fff;
    box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
    border-radius: 5px;
    border-top: 5px solid #ff6600; /* Warna yang sama dengan tombol Verify Email */
}

.footer {
    text-align: center;
    margin-top: 20px;
    color: #666;
    padding-top: 10px;
    border-top: 1px solid #ddd;
}
.btn-verify-email {
    background-color: #ff6600;
    color: #fff;
    padding: 12px 24px;
    border-radius: 5px;
    text-decoration: none;
    display: block;
    text-align: center;
    margin: 20px auto;
    transition: background-color 0.3s ease;
}

.btn-verify-email:hover {
    background-color: #ff3300;
}
.message {
    background-color: #f9f9f9;
    padding: 20px;
    border: 1px solid #ddd;
    border-radius: 5px;
    margin-top: 20px;
}

p {
    font-size: 18px;
    margin-top: 15px;
    line-height: 1.6;
}

h1 {
    text-align: center;
    color: #333;
    border-bottom: 2px solid #ff6600;
    padding-bottom: 10px;
}

        </style>
    </head>
    <body>
        <div class="container">
            <h1>Website Pendaftaran SMAN 1 Bunguran Barat</h1>
            <div class="message">
                <p>Hello, <strong>` + name + `</strong>,</p>
                <p>Terimakasih sudah register di website ini. Silahkan tekan tombol verify email untuk melanjutkan</p>
                <p><strong>Support Team:</strong> <a href="mailto:laodesaady12345@gmail.com">laodesaady12345@gmail.com</a></p>
                <a href="` + verificationLink + `" class="btn btn-verify-email">Verify Email</a>
            </div>
            <div class="footer">
                <p>&copy; 2023 laode saady. All rights reserved.</p>
            </div>
        </div>
    </body>
    </html>
    `

	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", sender)
	m.SetHeader("To", recipient)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", emailBody)

	d := gomail.NewDialer(smtpServer, smtpPort, smtpUsername, smtpPassword)

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
