package mail

import (
	"fmt"
	"net/smtp"
)

/// Функция SendEmail создает письмо для подтвержения регистрации \\\

func SendEmail(from, password, addressee, name, surname, confirmCode string) error {
	/// Создание тела сообщения \\\
	subject := "Confirmation of registration"
	body := fmt.Sprintf("Hello, %s %s. To register successfully, you need to confirm your mail.\n\nYour confirmation code is: %s\n\nBest regards,Your App", name, surname, confirmCode)
	msg := "From: " + from + "\n" +
		"To: " + addressee + "\n" +
		"Subject: " + subject + "\n\n" +
		body
	/// Выбор почтового сервиса \\\
	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, password, "smtp.gmail.com"),
		from, []string{addressee}, []byte(msg))
	if err != nil {
		return err
	}
	return nil
}

/// Функция SendAppealEmail создает письмо уведомления для обратной связи \\\

func SendAppealEmail(from, password, addressee, fio, mailsubject string) error {
	/// Создание тела сообщения \\\
	body := fmt.Sprintf("Hello, %s. Thank you for contacting us. Your letter on the subject: '%s' has been received. It will be reviewed during the day. If you have not received an answer, then contact any messenger convenient for you in the 'Contacts' section.\n\nBest regards,Your App", fio, mailsubject)
	msg := "From: " + from + "\n" +
		"To: " + addressee + "\n" +
		"Subject: " + mailsubject + "\n\n" +
		body
	/// Выбор почтового сервиса \\\
	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, password, "smtp.gmail.com"),
		from, []string{addressee}, []byte(msg))
	if err != nil {
		return err
	}
	fmt.Print(msg)
	return nil
}
