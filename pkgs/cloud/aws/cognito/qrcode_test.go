package cognito

import "testing"

func TestQrCodeGeneration(t *testing.T) {
	_, err := generateQRCode("Hello there !", 256)
	if err != nil {
		return
	}
}
