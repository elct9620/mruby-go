package mruby_test

import (
	"fmt"
	"testing"

	"github.com/cucumber/godog"
	"github.com/elct9620/mruby-go"
	"github.com/google/go-cmp/cmp"
)

const SuiteSuccessCode = 0

type RubyFeature struct {
	mrb *mruby.State
	ret mruby.Value
}

func (feat *RubyFeature) iExecuteRubyCode(code *godog.DocString) (err error) {
	feat.ret, err = feat.mrb.LoadString(code.Content)
	return err
}

func (feat *RubyFeature) thereShouldReturnInteger(expected int) error {
	actual, ok := feat.ret.(int)
	if !ok {
		return fmt.Errorf("expected integer, got %T", feat.ret)
	}

	if actual != expected {
		return fmt.Errorf("expected %d, got %d", expected, actual)
	}

	return nil
}

func (feat *RubyFeature) thereShouldReturnTrue() error {
	actual, ok := feat.ret.(bool)
	if !ok {
		return fmt.Errorf("expected bool, got %T", feat.ret)
	}

	if !actual {
		return fmt.Errorf("expected true, got false")
	}

	return nil
}

func (feat *RubyFeature) thereShouldReturnFalse() error {
	actual, ok := feat.ret.(bool)
	if !ok {
		return fmt.Errorf("expected bool, got %T", feat.ret)
	}

	if actual {
		return fmt.Errorf("expected false, got true")
	}

	return nil
}

func (feat *RubyFeature) thereShouldReturnNil() error {
	if feat.ret != nil {
		return fmt.Errorf("expected nil, got %T", feat.ret)
	}

	return nil
}

func (feat *RubyFeature) thereShouldReturnString(expected string) error {
	actual, ok := feat.ret.(string)
	if !ok {
		return fmt.Errorf("expected string, got %T", feat.ret)
	}

	if !cmp.Equal(actual, expected) {
		return fmt.Errorf("string not matched %s", cmp.Diff(actual, expected))
	}

	return nil
}

func (feat *RubyFeature) thereShouldReturnSymbol(expected string) error {
	ret, ok := feat.ret.(mruby.Symbol)
	if !ok {
		return fmt.Errorf("expected symbol, got %T", feat.ret)
	}

	actual := feat.mrb.SymbolName(ret)
	if actual != expected {
		return fmt.Errorf("expected %s, got %s", expected, actual)
	}

	return nil
}

func (feat *RubyFeature) thereShouldReturnObject() error {
	_, ok := feat.ret.(mruby.RBasic)
	if !ok {
		return fmt.Errorf("expected object, got %T", feat.ret)
	}

	return nil
}

func InitializeScenario(s *godog.ScenarioContext) {
	feat := RubyFeature{
		mrb: mruby.New(),
	}

	s.Step(`^I execute ruby code:$`, feat.iExecuteRubyCode)
	s.Step(`^there should return integer (-?\d+)$`, feat.thereShouldReturnInteger)
	s.Step(`^there should return true$`, feat.thereShouldReturnTrue)
	s.Step(`^there should return false$`, feat.thereShouldReturnFalse)
	s.Step(`^there should return nil$`, feat.thereShouldReturnNil)
	s.Step(`^there should return string "([^"]*)"$`, feat.thereShouldReturnString)
	s.Step(`^there should return symbol "([^"]*)"$`, feat.thereShouldReturnSymbol)
	s.Step(`^there should return object$`, feat.thereShouldReturnObject)
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t,
		},
	}

	if suite.Run() != SuiteSuccessCode {
		t.Fatal("Non-zero exit code")
	}
}
