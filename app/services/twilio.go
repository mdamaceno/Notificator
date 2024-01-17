package services

import (
	"log"
	"os"

	twilioClient "github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

type TwilioSMSService struct{}
type TwilioWhatsappService struct{}

func (s TwilioSMSService) Send(receivers []string, message string) []error {
	var errList []error

	client := twilioClient.NewRestClientWithParams(twilioClient.ClientParams{
		Username: os.Getenv("SMS_USERNAME"),
		Password: os.Getenv("SMS_PASSWORD"),
	})
	params := &twilioApi.CreateMessageParams{}
	params.SetBody(message)
	params.SetFrom(os.Getenv("SMS_FROM_NUMBER"))

	for _, receiver := range receivers {
		params.SetTo(receiver)
		_, err := client.Api.CreateMessage(params)
		if err != nil {
			log.Println(err)
			errList = append(errList, err)
		}
	}

	return errList
}

func (s TwilioWhatsappService) Send(receivers []string, message string) []error {
	var errList []error

	client := twilioClient.NewRestClientWithParams(twilioClient.ClientParams{
		Username: os.Getenv("WHATSAPP_USERNAME"),
		Password: os.Getenv("WHATSAPP_PASSWORD"),
	})

	params := &twilioApi.CreateMessageParams{}
	params.SetBody(message)
	params.SetFrom("whatsapp:" + os.Getenv("WHATSAPP_FROM_NUMBER"))

	for _, receiver := range receivers {
		params.SetTo("whatsapp:" + receiver)
		res, err := client.Api.CreateMessage(params)
		if err != nil {
			if *res.ErrorCode == 21608 && numberFromBrazil(receiver) {
				receiver = s.remove9DigitBrazil(receiver)
				params.SetTo("whatsapp:" + receiver)
				res, err = client.Api.CreateMessage(params)
				if err != nil {
					log.Println(err)
					errList = append(errList, err)
				}
			} else {
				log.Println(err)
				errList = append(errList, err)
			}
		}
	}

	return errList
}

func (s TwilioWhatsappService) remove9DigitBrazil(number string) string {
	countryCode := number[0:3]
	phoneNumber := number[3:]

	if len(phoneNumber) == 12 && countryCode == "+55" {
		return countryCode + phoneNumber[1:]
	}

	return number
}

func numberFromBrazil(number string) bool {
	countryCode := number[0:3]

	return countryCode == "+55"
}
