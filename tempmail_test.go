package gotempmail

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestTempMailNew(t *testing.T) {
	tempmail := New().
		Address("Testing123@test.com").
		Password("testing123")
	if tempmail == nil {
		t.Error("tempmail is nil")
	}
}

func TestTempMailValidate(t *testing.T) {
	var tempmail TempMail
	err := tempmail.Validate()
	if err == nil {
		t.Errorf("validation should have failed %s", err)
	}

	tempmail.password = "wee wee"
	err = tempmail.Validate()
	if err == nil {
		t.Errorf("validation should have failed %s", err)
	}

	tempmail.Email = "a@b.com"
	err = tempmail.Validate()
	if err != nil {
		t.Errorf("validation should not have failed %s", err)
	}

	tempmail.Email = "com"
	err = tempmail.Validate()
	if err == nil {
		t.Errorf("validation should have failed %s", err)
	}
}

const (
	MX_TOOLBOX_SEND_URL      = "https://mxtoolbox.com/public/tools/EmailHeaders.aspx"
	MX_TOOLBOX_EXPECTED_CODE = http.StatusOK
	MX_TOOLBOX_DATA_FIELD    = "ctl00$ContentPlaceHolder1$txtEmail"
	TEST_WAIT_TIME           = time.Second * 5
	TEST_MAX_TRIES           = 10
)

func TestTempMail(t *testing.T) {
	// Get domains
	domains, err := GetDomains()
	if err != nil {
		t.Error(err)
	}

	if domains == nil {
		t.Error("Nil Domains")
	}
	if len(domains) == 0 {
		t.Error("No domains")
	}

	// Create the email
	tempmail, err := New().
		Address("testing" + fmt.Sprintf("%d",
			time.Now().Unix()) + "@" + domains[0]).
		Password("password123").
		Build()
	if err != nil {
		t.Errorf("tempmail err is %s", err)
	}

	emails, err := tempmail.GetEmails()
	if err != nil {
		t.Errorf("tempmail GetEmails () err %s", err)
	}

	if emails == nil {
		t.Errorf("emails are nil")
	}

	if len(emails) != 0 {
		t.Errorf("there should not be any emails yet")
	}

	emails, err = tempmail.GetEmails()
	if err != nil {
		t.Errorf("tempmail GetEmails () err %s", err)
	}

	for _, email := range emails {
		details, err := tempmail.GetEmailDetails(email)
		if err != nil {
			t.Errorf("error getting email details %s", err)
		}

		if details.Subject != email.Subject {
			t.Errorf("wrong subject")
		}
	}
}

func TestEmailUnmarshal(t *testing.T) {
	testData := `{
  "hydra:member": [
    {
      "@id": "string",
      "@type": "string",
      "@context": "string",
      "id": "string",
      "accountId": "string",
      "msgid": "string",
      "from": {
        "address": "from@example.com",
        "name": "John Doe"
      },
      "to": [
        {
          "address": "receiver@example.com",
          "name": "John Doe"
        }
      ],
      "subject": "string",
      "intro": "string",
      "seen": true,
      "isDeleted": true,
      "hasAttachments": true,
      "size": 0,
      "downloadUrl": "string",
      "createdAt": "2023-06-04T15:36:13.408Z",
      "updatedAt": "2023-06-04T15:36:13.408Z"
    }
  ],
  "hydra:totalItems": 0,
  "hydra:view": {
    "@id": "string",
    "@type": "string",
    "hydra:first": "string",
    "hydra:last": "string",
    "hydra:previous": "string",
    "hydra:next": "string"
  },
  "hydra:search": {
    "@type": "string",
    "hydra:template": "string",
    "hydra:variableRepresentation": "string",
    "hydra:mapping": [
      {
        "@type": "string",
        "variable": "string",
        "property": "string",
        "required": true
      }
    ]
  }
}`

	var emails emailsJson
	err := json.Unmarshal([]byte(testData), &emails)
	if err != nil {
		t.Errorf("Cannot unmarshal %s", err)
	}
}

func TestEmailMessageUnmarshal(t *testing.T) {
	testData := `{
  "@context": "string",
  "@id": "string",
  "@type": "string",
  "id": "string",
  "accountId": "string",
  "msgid": "string",
  "from": {
    "address": "from@example.com",
    "name": "John Doe"
  },
  "to": [
    {
      "address": "receiver@example.com",
      "name": "John Doe"
    }
  ],
  "cc": [
    {
      "address": "cc@example.com",
      "name": "John Doe"
    }
  ],
  "bcc": [
    {
      "address": "bcc@example.com",
      "name": "John Doe"
    }
  ],
  "subject": "string",
  "seen": true,
  "flagged": true,
  "isDeleted": true,
  "verifications": [
    "string"
  ],
  "retention": true,
  "retentionDate": "2023-06-04T15:39:16.523Z",
  "text": "string",
  "html": [
    "string"
  ],
  "hasAttachments": true,
  "attachments": [
    {
      "id": "ATTACH000001",
      "filename": "happy.png",
      "contentType": "image/png",
      "disposition": "attachment",
      "transferEncoding": "base64",
      "related": false,
      "size": 666,
      "downloadUrl": "/messages/id/attachment/ATTACH000001"
    }
  ],
  "size": 0,
  "downloadUrl": "string",
  "createdAt": "2023-06-04T15:39:16.523Z",
  "updatedAt": "2023-06-04T15:39:16.523Z"
}`

	var message EmailDetails
	err := json.Unmarshal([]byte(testData), &message)
	if err != nil {
		t.Errorf("Cannot unmarshal %s", err)
	}
}

func TestEmailAddressUnmarshal(t *testing.T) {
	testData := `{
      "address": "receiver@example.com",
      "name": "John Doe"
  }`

	var addr EmailAddr
	err := json.Unmarshal([]byte(testData), &addr)
	if err != nil {
		t.Errorf("Cannot unmarshal %s", err)
	}
}
