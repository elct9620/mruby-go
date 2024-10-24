package mruby_test

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/cucumber/godog"
	"github.com/elct9620/mruby-go"
	"github.com/google/go-cmp/cmp"
)

var opts = godog.Options{
	Tags:   "~@wip",
	Format: "pretty",
	Paths:  []string{"features"},
}

const SuiteSuccessCode = 0

type RubyFeature struct {
	mrb *mruby.State
	ret mruby.Value
	exc mruby.RException
}

func (feat *RubyFeature) iExecuteRubyCode(code *godog.DocString) error {
	ret, err := feat.mrb.LoadString(code.Content)
	if err != nil {
		exc, ok := err.(mruby.RException)
		if !ok {
			return err
		}
		feat.exc = exc
	}
	feat.ret = ret

	return nil
}

func (feat *RubyFeature) thereShouldReturnInteger(expected int) error {
	actual, ok := feat.ret.(int)
	if !ok {
		return fmt.Errorf("expected integer, got %T (%+v)", feat.ret, feat.ret)
	}

	if actual != expected {
		return fmt.Errorf("expected %d, got %d", expected, actual)
	}

	return nil
}

func (feat *RubyFeature) thereShouldReturnTrue() error {
	actual, ok := feat.ret.(bool)
	if !ok {
		return fmt.Errorf("expected bool, got %T (%+v)", feat.ret, feat.ret)
	}

	if !actual {
		return fmt.Errorf("expected true, got false")
	}

	return nil
}

func (feat *RubyFeature) thereShouldReturnFalse() error {
	actual, ok := feat.ret.(bool)
	if !ok {
		return fmt.Errorf("expected bool, got %T (%+v)", feat.ret, feat.ret)
	}

	if actual {
		return fmt.Errorf("expected false, got true")
	}

	return nil
}

func (feat *RubyFeature) thereShouldReturnNil() error {
	if feat.ret != nil {
		return fmt.Errorf("expected nil, got %T (%+v)", feat.ret, feat.ret)
	}

	return nil
}

func (feat *RubyFeature) thereShouldReturnString(expected string) error {
	actual, ok := feat.ret.(string)
	if !ok {
		return fmt.Errorf("expected string, got %T (%+v)", feat.ret, feat.ret)
	}

	if !cmp.Equal(actual, expected) {
		return fmt.Errorf("string not matched %s", cmp.Diff(actual, expected))
	}

	return nil
}

func (feat *RubyFeature) thereShouldReturnStringLike(expected string) error {
	actual, ok := feat.ret.(string)
	if !ok {
		return fmt.Errorf("expected string, got %T (%+v)", feat.ret, feat.ret)
	}

	expr, err := regexp.Compile(expected)
	if err != nil {
		return err
	}

	if !expr.MatchString(actual) {
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
	_, ok := feat.ret.(mruby.RObject)
	if !ok {
		return fmt.Errorf("expected object, got %T", feat.ret)
	}

	return nil
}

func (feat *RubyFeature) thereShouldReturnClass(expected string) error {
	class, ok := feat.ret.(*mruby.Class)

	if !ok {
		return fmt.Errorf("expected class, got %T", feat.ret)
	}

	actual := feat.mrb.ClassName(class)
	if actual != expected {
		return fmt.Errorf("expected %s, got %s", expected, actual)
	}

	return nil
}

func (feat *RubyFeature) thereShouldReturnModule(expected string) error {
	module, ok := feat.ret.(*mruby.Module)

	if !ok {
		return fmt.Errorf("expected module, got %T", feat.ret)
	}

	actual := feat.mrb.ClassName(module)
	if actual != expected {
		return fmt.Errorf("expected %s, got %s", expected, actual)
	}

	return nil
}

func (feat *RubyFeature) thereShouldReturnAnArray(doc *godog.DocString) error {
	actual, ok := feat.ret.([]mruby.Value)

	if !ok {
		return fmt.Errorf("expected array, got %T", feat.ret)
	}

	actualStr := fmt.Sprintf("%+v", actual)
	expectedStr := doc.Content

	if !cmp.Equal(actualStr, expectedStr) {
		return fmt.Errorf("array not matched %s", cmp.Diff(actualStr, expectedStr))

	}

	return nil
}

func (feat *RubyFeature) thereShouldReturnAHash(doc *godog.DocString) error {
	actual, ok := feat.ret.(map[mruby.Value]mruby.Value)

	if !ok {
		return fmt.Errorf("expected hash, got %T", feat.ret)
	}

	actualStr := fmt.Sprintf("%+v", actual)
	expectedStr := doc.Content

	if !cmp.Equal(actualStr, expectedStr) {
		return fmt.Errorf("hash not matched %s", cmp.Diff(actualStr, expectedStr))

	}

	return nil
}

func (feat *RubyFeature) theExceptionMessageShouldBe(expected string) error {
	actual := feat.exc.Error()
	if actual != expected {
		return fmt.Errorf("expected %s, got %s", expected, actual)
	}

	return nil
}

func InitializeScenario(s *godog.ScenarioContext) {
	mrb, err := mruby.New()
	if err != nil {
		panic(err)
	}

	feat := RubyFeature{
		mrb: mrb,
	}

	s.Step(`^I execute ruby code:$`, feat.iExecuteRubyCode)
	s.Step(`^there should return integer (-?\d+)$`, feat.thereShouldReturnInteger)
	s.Step(`^there should return true$`, feat.thereShouldReturnTrue)
	s.Step(`^there should return false$`, feat.thereShouldReturnFalse)
	s.Step(`^there should return nil$`, feat.thereShouldReturnNil)
	s.Step(`^there should return string "([^"]*)"$`, feat.thereShouldReturnString)
	s.Step(`^there should return string like "([^"]*)"$`, feat.thereShouldReturnStringLike)
	s.Step(`^there should return symbol "([^"]*)"$`, feat.thereShouldReturnSymbol)
	s.Step(`^there should return object$`, feat.thereShouldReturnObject)
	s.Step(`^there should return class "([^"]*)"$`, feat.thereShouldReturnClass)
	s.Step(`^there should return module "([^"]*)"$`, feat.thereShouldReturnModule)
	s.Step(`^there should return an array$`, feat.thereShouldReturnAnArray)
	s.Step(`^there should return a hash$`, feat.thereShouldReturnAHash)
	s.Step(`^the exception message should be "([^"]*)"$`, feat.theExceptionMessageShouldBe)
}

func init() {
	godog.BindFlags("godog.", flag.CommandLine, &opts)
}

func TestFeatures(t *testing.T) {
	o := opts
	o.TestingT = t

	if o.Format == "junit" {
		output, err := os.OpenFile("junit.xml", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			t.Fatal(err)
		}
		defer output.Close()
		o.Output = output
	}

	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options:             &o,
	}

	if suite.Run() != SuiteSuccessCode {
		t.Fatal("Non-zero exit code")
	}
}
