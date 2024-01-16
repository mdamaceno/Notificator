package twilio

import (
	"log"
	"os"

	twilioClient "github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

type TwilioSMSService struct{}

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
