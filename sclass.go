package mruby

var _ RClass = &SingletonClass{}

type SingletonClass struct {
	super RClass
	class RClass
	flags uint32
	mt    methodTable
	iv    ivTable
}

func (c *SingletonClass) Class() RClass {
	return c.class
}

func (c *SingletonClass) Flags() uint32 {
	return c.flags
}

func (c *SingletonClass) ivPut(sym Symbol, val Value) {
	if c.iv == nil {
		c.iv = make(ivTable)
	}

	c.iv[sym] = val
}

func (c *SingletonClass) ivGet(sym Symbol) Value {
	if c.iv == nil {
		return nil
	}

	return c.iv[sym]
}

func (c *SingletonClass) mtPut(sym Symbol, method Method) {
	if c.mt == nil {
		c.mt = make(methodTable)
	}

	c.mt[sym] = method
}

func (c *SingletonClass) mtGet(sym Symbol) Method {
	if c.mt == nil {
		return nil
	}

	return c.mt[sym]
}

func (c *SingletonClass) Super() RClass {
	return c.super
}
