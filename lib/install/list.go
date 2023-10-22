package install

func list() *installList {
	list := installList{
		{"nvm", nil, true, nil},
		{"gcc", &pkg{"gcc", "gcc", ""}, false, nil},
		{"gh", &pkg{"gh", "gh", ""}, false, nil},
		{"g++", &pkg{"gcc-c++", "gcc-c++", ""}, false, nil},
		{"go", nil, true, nil},
		{"flatpak", &pkg{"flatpak", "flatpak", ""}, false, ConfigureFlatpak},
		{"codium", &pkg{"", "", "com.vscodium.codium"}, false, nil},
		{"starship", nil, true, nil},
	}
	return &list
}
