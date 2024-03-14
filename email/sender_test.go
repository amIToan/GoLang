package email

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"sgithub.com/techschool/simplebank/util"
)

func TestSendEmailWithGmail(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	config, err := util.LoadConfigDB_Server("..")
	require.NoError(t, err)
	fmt.Println("username", config.EmailSenderName)
	fmt.Println("username", config.EmailSenderAddress)
	sender := CreateNewEmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)

	subject := "A test email"
	content := `
	<h1>Hello world</h1>
	<p>This is a test message from <a href="http://techschool.guru">Tech School</a></p>
	`
	to := []string{"djitoan@gmail.com"}
	attachFiles := []string{"../README.md"}

	err = sender.SendEmail(subject, content, to, nil, nil, attachFiles)
	require.NoError(t, err)
}
