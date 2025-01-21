package handler

import "gopkg.in/gomail.v2"

var Dialer *gomail.Dialer

func init() {
	Dialer = gomail.NewDialer("smtp.gmail.com", 587, "sanek.tursumetov@gmail.com", "stlv hite ixhw kbed")
}
