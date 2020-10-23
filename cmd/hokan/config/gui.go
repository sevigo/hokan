package config

func defaultGUI(c *Config) {
	if c.GUI.AppName == "" {
		c.GUI.AppName = "Hokan UI"
	}
	c.GUI.BaseDir = "../../resources/"
	c.GUI.DevTools = false
}
