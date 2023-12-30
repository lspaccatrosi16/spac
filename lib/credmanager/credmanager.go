package credmanager

import "github.com/lspaccatrosi16/go-cli-tools/credential"

func Credential() error {
	return credential.StandaloneManager()
}
