package mail

import (
	"testing"

	"github.com/HamedBlue1381/hamed-bank/util"
	"github.com/stretchr/testify/require"
)

func TestSendEmailWithGmail(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	config, err := util.LoadConfig("..")
	require.NoError(t, err)

	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)
	//TODO: add generator and content's constant var
	subject := "Cute Email :)"
	content := MessageGenerator()
	to := []string{ToEmailAddress}
	attachFiles := []string{"../Email.txt"}

	err = sender.SendEmail(subject, content, to, nil, nil, attachFiles)
	require.NoError(t, err)
}
