package creature

func ReduceNeed(c *Creature, needType NeedType, amount float64) {
	for i, n := range c.Needs {
		if n.Type == needType {
			c.Needs[i].Value -= amount
			if c.Needs[i].Value < 0 {
				c.Needs[i].Value = 0
			}
			break
		}
	}
}