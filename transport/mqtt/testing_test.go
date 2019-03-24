package mqtt

import "os"

const integrationDisabledMsg = "Integration tests disabled"

var enableIntegrationtests = false

func init() {
	if os.Getenv("BIBLIOTEK_TEST_INTEGRATION") == "1" {
		enableIntegrationtests = true
	}
}
