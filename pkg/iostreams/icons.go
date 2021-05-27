package iostreams

func (c *ColorPalette) GreenCheck() string {
	return c.Green("✓")
}

func (c *ColorPalette) FailedIcon() string {
	return c.Red("x")
}

func (c *ColorPalette) WarnIcon() string {
	return c.Yellow("!")
}

func (c *ColorPalette) RedCheck() string {
	return c.Red("✓")
}

func (c *ColorPalette) ProgressIcon() string {
	return c.Blue("•")
}

func (c *ColorPalette) DotWarnIcon() string {
	return c.Yellow("•")
}
