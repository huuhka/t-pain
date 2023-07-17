package openai

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

func LoginWithDefaultCredential() (*azidentity.DefaultAzureCredential, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, fmt.Errorf("unable to get default credential: %w", err)
	}

	return cred, nil
}