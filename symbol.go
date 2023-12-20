package mruby

type Symbol = uint32

func (c *State) Intern(str string) Symbol {
	if sym, ok := c.symbolTable[str]; ok {
		return sym
	}

	c.symbolIndex++
	c.symbolTable[str] = Symbol(c.symbolIndex)

	return Symbol(c.symbolIndex)
}

func (c *State) SymbolName(sym Symbol) string {
	for k, v := range c.symbolTable {
		if v == sym {
			return k
		}
	}

	return ""
}
