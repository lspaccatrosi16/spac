package config

func list() configList {
	return configList{
		{"starship", "starship.toml", "$HOME/.config"},
	}
}
