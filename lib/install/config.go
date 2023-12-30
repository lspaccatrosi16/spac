package install

import "github.com/lspaccatrosi16/spac/lib/install/sudo"

func ConfigureFlatpak() error {
	return sudo.RunSudo("flatpak remote-add --if-not-exists flathub https://dl.flathub.org/repo/flathub.flatpakrepo")
}
