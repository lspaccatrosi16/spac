package install

func list() *installList {
	list := installList{
		{"nvm", "NodeJS version manager", nil, true, nil},
		{"gcc", "C/C++ compiler", &pkg{"gcc", nil, "gcc", ""}, false, nil},
		{"gh", "Github CLI", &pkg{"gh", nil, "gh", ""}, false, nil},
		{"g++", "C/C++ compiler", &pkg{"gcc-c++", nil, "gcc-c++", ""}, false, nil},
		{"go", "Golang compiler", nil, true, nil},
		{"flatpak", "Package manager", &pkg{"flatpak", nil, "flatpak", ""}, false, ConfigureFlatpak},
		{"codium", "Code editor", &pkg{"", nil, "", "com.vscodium.codium"}, false, nil},
		{"starship", "Shell prompt", nil, true, nil},
		{"aup", "Package manager", nil, true, nil},
		{"scaffold", "Project templating tool", &pkg{"", &aupPkg{"lspaccatrosi16/scaffold", "scaffold-linux", "scaffold"}, "", ""}, false, nil},
		{"releasetool", "Automatic github releaser", &pkg{"", &aupPkg{"lspaccatrosi16/releasetool", "releasetool-linux", "release"}, "", ""}, false, nil},
		{"bupload", "Quick uploader to / from cloud buckets", &pkg{"", &aupPkg{"lspaccatrosi16/bupload", "bupload-linux", "bupload"}, "", ""}, false, nil},
		{"motionTui", "TUI for Motion", &pkg{"", &aupPkg{"lspaccatrosi16/motion-tui", "motion-linux", "motion"}, "", ""}, false, nil},
		{"ts-go", "Converts Typescript types into Golang types", &pkg{"", &aupPkg{"lspaccatrosi16/ts-go", "ts-go-linux", "ts-go"}, "", ""}, false, nil},
	}
	return &list
}
