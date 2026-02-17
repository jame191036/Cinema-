package services

import (
	"fmt"
	"log"
	"net/smtp"
	"strings"
)

type EmailService struct {
	host     string
	port     string
	username string
	password string
	from     string
}

func NewEmailService(host, port, username, password, from string) *EmailService {
	return &EmailService{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
	}
}

func (s *EmailService) IsConfigured() bool {
	return s.host != "" && s.username != "" && s.password != ""
}

type BookingConfirmationData struct {
	UserName   string
	UserEmail  string
	BookingID  string
	Seats      []string
	ShowtimeID string
	OccurredAt string
}

func (s *EmailService) SendBookingConfirmation(data BookingConfirmationData) error {
	if !s.IsConfigured() {
		log.Printf("Email service not configured, skipping email to %s", data.UserEmail)
		return nil
	}

	subject := "🎬 ยืนยันการจองตั๋วเรียบร้อยแล้ว!"
	body := buildEmailBody(data)

	msg := buildMIMEMessage(s.from, data.UserEmail, subject, body)

	addr := fmt.Sprintf("%s:%s", s.host, s.port)
	auth := smtp.PlainAuth("", s.username, s.password, s.host)

	err := smtp.SendMail(addr, auth, s.from, []string{data.UserEmail}, []byte(msg))
	if err != nil {
		return fmt.Errorf("failed to send email to %s: %w", data.UserEmail, err)
	}

	log.Printf("Booking confirmation email sent to %s (bookingId: %s)", data.UserEmail, data.BookingID)
	return nil
}

func buildEmailBody(data BookingConfirmationData) string {
	seats := strings.Join(data.Seats, ", ")
	return fmt.Sprintf(`
สวัสดีคุณ %s,

การจองตั๋วภาพยนตร์ของคุณได้รับการยืนยันเรียบร้อยแล้ว! 🎉

รายละเอียดการจอง:
━━━━━━━━━━━━━━━━━━━━━━━━
  Booking ID : %s
  ที่นั่ง     : %s
  วันที่จอง  : %s
━━━━━━━━━━━━━━━━━━━━━━━━

กรุณาแสดง Booking ID ณ จุดรับบัตรก่อนเข้าฉาย

ขอบคุณที่ใช้บริการ 🎬
Cinema Booking System
`, data.UserName, data.BookingID, seats, data.OccurredAt)
}

func buildMIMEMessage(from, to, subject, body string) string {
	header := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n",
		from, to, subject,
	)
	return header + body
}
