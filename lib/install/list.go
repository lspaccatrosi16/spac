package install

func list() *installList {
	list := installList{
		{"nvm", nil, true, nil},
		{"gcc", &pkg{"gcc", nil, "gcc", ""}, false, nil},
		{"gh", &pkg{"gh", nil, "gh", ""}, false, nil},
		{"g++", &pkg{"gcc-c++", nil, "gcc-c++", ""}, false, nil},
		{"go", nil, true, nil},
		{"flatpak", &pkg{"flatpak", nil, "flatpak", ""}, false, ConfigureFlatpak},
		{"codium", &pkg{"", nil, "", "com.vscodium.codium"}, false, nil},
		{"starship", nil, true, nil},
		{"aup", nil, true, nil},
		{"scaffold", &pkg{"", &aupPkg{"lspaccatrosi16/scaffold", "scaffold-linux", "scaffold"}, "", ""}, false, nil},
		{"releasetool", &pkg{"", &aupPkg{"lspaccatrosi16/releasetool", "releasetool-linux", "release"}, "", ""}, false, nil},
	}
	return &list
}
