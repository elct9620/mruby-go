mruby-go
===

[![Go Reference](https://pkg.go.dev/badge/github.com/elct9620/mruby-go.svg)](https://pkg.go.dev/github.com/elct9620/mruby-go)
[![Test](https://github.com/elct9620/mruby-go/actions/workflows/test.yml/badge.svg)](https://github.com/elct9620/mruby-go/actions/workflows/test.yml)
[![Maintainability](https://api.codeclimate.com/v1/badges/62c60dab046a3d550e78/maintainability)](https://codeclimate.com/github/elct9620/mruby-go/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/62c60dab046a3d550e78/test_coverage)](https://codeclimate.com/github/elct9620/mruby-go/test_coverage)

The pure go mruby virtual machine implementation.

## Roadmap

The priority task is to make the virtual machine available with limited functions, it still depends on `mrbc` command to compile RiteBinary.

## MRB_API

Golang has public and private method design and we can attach method to a struct. Therefore all public method is attach to `*mruby.State` in mruby-go as preferred method.

```go
func (mrb *State) ObjectInstanceVariableGet(obj RObject, name Symbol) Value {
  return obj.ivGet(name)
}

// Prefer
func (mrb *State) ClassName(class RClass) string {
	if class == nil {
		return ""
	}

	name := mrb.ObjectInstanceVariableGet(class, _classname(mrb)) // <- Prefer this
	if name == nil {
		return ""
	}

	return name.(string)
}

// Avoid
func (mrb *State) ClassName(class RClass) string {
	if class == nil {
		return ""
	}

	name := class.ivGet(_classname(mrb)) // <- Avoid this
	if name == nil {
		return ""
	}

	return name.(string)
}
```
