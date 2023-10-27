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
		{"starship", "Terminal host", nil, true, nil},
		{"aup", "Package manager", nil, true, nil},
		{"scaffold", "Project templating tool", &pkg{"", &aupPkg{"lspaccatrosi16/scaffold", "scaffold-linux", "scaffold"}, "", ""}, false, nil},
		{"releasetool", "Automatic github releaser", &pkg{"", &aupPkg{"lspaccatrosi16/releasetool", "releasetool-linux", "release"}, "", ""}, false, nil},
	}
	return &list
}
