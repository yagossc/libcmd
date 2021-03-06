package libcmd_test

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/ibraimgm/libcmd"
)

func compareHelpOutput(app *libcmd.App, args []string, goldenfile string) error {
	bytes, err := ioutil.ReadFile(goldenfile)
	if err != nil {
		return err
	}
	expected := string(bytes)

	var b strings.Builder
	app.Options.HelpOutput = &b

	if len(args) == 0 {
		app.Help()
	} else if err := app.ParseArgs(args); err != nil {
		return err
	}

	actual := b.String()

	if expected != actual {
		return fmt.Errorf("Wrong output. Expected:\n>>>\n%s\n<<<\nActual:\n>>>\n%s\n<<<", expected, actual)
	}

	return nil
}

func TestBasic(t *testing.T) {
	app := libcmd.NewApp("app", "some brief description")

	if err := compareHelpOutput(app, []string{}, "testdata/basic.golden"); err != nil {
		t.Error(err)
	}
}

func TestNoBrief(t *testing.T) {
	app := libcmd.NewApp("app", "")

	if err := compareHelpOutput(app, []string{}, "testdata/nobrief.golden"); err != nil {
		t.Error(err)
	}
}

func TestNoUsage(t *testing.T) {
	app := libcmd.NewApp("app", "some brief description")
	app.Usage = "-"

	if err := compareHelpOutput(app, []string{}, "testdata/nousage.golden"); err != nil {
		t.Error(err)
	}
}

func TestCustomUsage(t *testing.T) {
	app := libcmd.NewApp("app", "some brief description")
	app.Usage = "my custom usage text"

	if err := compareHelpOutput(app, []string{}, "testdata/usage.golden"); err != nil {
		t.Error(err)
	}
}

func TestLong(t *testing.T) {
	app := libcmd.NewApp("app", "some brief description")
	app.Long = "this is a very long description"

	if err := compareHelpOutput(app, []string{}, "testdata/long.golden"); err != nil {
		t.Error(err)
	}
}

func TestOperands(t *testing.T) {
	app := libcmd.NewApp("app", "some brief description")
	app.AddOperand("src", "")
	app.AddOperand("dst", "")

	if err := compareHelpOutput(app, []string{}, "testdata/operands.golden"); err != nil {
		t.Error(err)
	}
}

func TestOperandsMod(t *testing.T) {
	app := libcmd.NewApp("app", "some brief description")
	app.AddOperand("src", "?")
	app.AddOperand("dst", "*")

	if err := compareHelpOutput(app, []string{}, "testdata/operands-mod.golden"); err != nil {
		t.Error(err)
	}
}

func TestArgs(t *testing.T) {
	app := libcmd.NewApp("app", "some brief description")
	app.Long = "this is a very long description"

	app.String("astring", 's', "", "Sets a string value.")
	app.Int("aint", 'i', 0, "Sets a int value.")
	app.Int("", 'a', 0, "Sets the amount value.", "amount")
	app.Int("balance", -1, 0, "Sets the balance value.", "balance")

	if err := compareHelpOutput(app, []string{"-h"}, "testdata/args.golden"); err != nil {
		t.Error(err)
	}
}

func TestHelpMessageWithOperands(t *testing.T) {
	app := libcmd.NewApp("app", "some brief description")
	app.AddOperand("src", "?")
	app.AddOperand("dst", "*")

	if err := compareHelpOutput(app, []string{"-h", "test", "test2"}, "testdata/help-message-with-operands.golden"); err != nil {
		t.Error(err)
	}
}

func TestArgsPartial(t *testing.T) {
	app := libcmd.NewApp("app", "some brief description")
	app.Long = "this is a very long description"

	app.String("", 's', "", "Sets a string value.")
	app.Int("aint", 0, 0, "Sets a int value.")

	if err := compareHelpOutput(app, []string{"-h"}, "testdata/partial.golden"); err != nil {
		t.Error(err)
	}
}

func TestArgsNoHelp(t *testing.T) {
	app := libcmd.NewApp("app", "some brief description")
	app.Long = "this is a very long description"
	app.Options.SuppressHelpFlag = true

	app.String("astring", 's', "", "Sets a string value.")
	app.Int("aint", 'i', 0, "Sets a int value.")

	if err := compareHelpOutput(app, []string{}, "testdata/nohelp.golden"); err != nil {
		t.Error(err)
	}
}

func TestDefault(t *testing.T) {
	app := libcmd.NewApp("app", "some brief description")
	app.Long = "this is a very long description"

	app.String("astring", 's', "somevalue", "Sets a string value.")
	app.Int("aint", 'i', 100, "Sets a int value.")

	if err := compareHelpOutput(app, []string{"-h"}, "testdata/default.golden"); err != nil {
		t.Error(err)
	}
}

func TestCommand(t *testing.T) {
	app := libcmd.NewApp("app", "some brief description")
	app.Long = "this is a very long description"

	app.Command("add", "Sums two numbers.", nil)
	app.Command("sub", "Subtract two numbers.", nil)

	if err := compareHelpOutput(app, []string{"-h"}, "testdata/command.golden"); err != nil {
		t.Error(err)
	}
}

func TestCommandArgs(t *testing.T) {
	app := libcmd.NewApp("app", "some brief description")
	app.Long = "this is a very long description"

	app.String("astring", 's', "somevalue", "Sets a string value.")
	app.Int("aint", 'i', 100, "Sets a int value.")

	app.Command("add", "Sums two numbers.", nil)
	app.Command("sub", "Subtract two numbers.", nil)

	if err := compareHelpOutput(app, []string{"-h"}, "testdata/commandargs.golden"); err != nil {
		t.Error(err)
	}
}

func TestSubcommand(t *testing.T) {
	app := libcmd.NewApp("app", "some brief description")
	app.Long = "this is a very long description"

	app.String("astring", 's', "somevalue", "sets a string value")
	app.Int("aint", 'i', 100, "sets a int value")

	app.Command("add", "Sums two numbers.", func(cmd *libcmd.Cmd) {
		cmd.Long = "Runs a computation that returns the sum of two specified numbers."
	})
	app.Command("sub", "Subtract two numbers.", nil)

	if err := compareHelpOutput(app, []string{"add"}, "testdata/subcommand.golden"); err != nil {
		t.Error(err)
	}
}

func TestSubcommandShallow(t *testing.T) {
	app := libcmd.NewApp("app", "some brief description")
	app.Long = "this is a very long description"

	app.String("astring", 's', "somevalue", "sets a string value")
	app.Int("aint", 'i', 100, "sets a int value")

	app.Command("add", "Sums two numbers.", func(cmd *libcmd.Cmd) {
		cmd.Long = "Runs a computation that returns the sum of two specified numbers."
		cmd.Command("deep", "A deep subcommand.", nil)
	})
	app.Command("sub", "Subtract two numbers.", nil)

	if err := compareHelpOutput(app, []string{"add"}, "testdata/shallow.golden"); err != nil {
		t.Error(err)
	}
}

func TestSubcommandShallowOp(t *testing.T) {
	app := libcmd.NewApp("app", "some brief description")
	app.Long = "this is a very long description"

	app.String("astring", 's', "somevalue", "sets a string value")
	app.Int("aint", 'i', 100, "sets a int value")

	app.Command("add", "Sums two numbers.", func(cmd *libcmd.Cmd) {
		cmd.Long = "Runs a computation that returns the sum of two specified numbers."
		cmd.AddOperand("number1", "")
		cmd.AddOperand("number2", "")

		cmd.Command("deep", "A deep subcommand.", nil)
	})
	app.Command("sub", "Subtract two numbers.", nil)

	if err := compareHelpOutput(app, []string{"add"}, "testdata/shallow-op.golden"); err != nil {
		t.Error(err)
	}
}

func TestSubcommandShallowOpRepeat(t *testing.T) {
	app := libcmd.NewApp("app", "some brief description")
	app.Long = "this is a very long description"

	app.String("astring", 's', "somevalue", "sets a string value")
	app.Int("aint", 'i', 100, "sets a int value")

	app.Command("add", "Sums two numbers.", func(cmd *libcmd.Cmd) {
		cmd.Long = "Runs a computation that returns the sum of two specified numbers."
		cmd.AddOperand("number1", "")
		cmd.AddOperand("others", "*")

		cmd.Command("deep", "A deep subcommand.", nil)
	})
	app.Command("sub", "Subtract two numbers.", nil)

	if err := compareHelpOutput(app, []string{"add"}, "testdata/shallow-op-repeat.golden"); err != nil {
		t.Error(err)
	}
}

func TestSubcommandDeep(t *testing.T) {
	app := libcmd.NewApp("app", "some brief description")
	app.Long = "this is a very long description"

	app.String("astring", 's', "somevalue", "sets a string value")
	app.Int("aint", 'i', 100, "sets a int value")

	app.Command("add", "Sums two numbers.", func(cmd *libcmd.Cmd) {
		cmd.Long = "Runs a computation that returns the sum of two specified numbers."
		cmd.Command("deep", "A deep subcommand.", func(cmd *libcmd.Cmd) {
			cmd.Long = "This is a deep subcommand."
		})
	})
	app.Command("sub", "Subtract two numbers.", nil)

	if err := compareHelpOutput(app, []string{"add", "deep"}, "testdata/deep.golden"); err != nil {
		t.Error(err)
	}
}

func TestHelpChoice(t *testing.T) {
	tests := []struct {
		choices  []string
		defValue string
		file     string
	}{
		{choices: []string{"foo", "bar", "baz"}, file: "testdata/choice-nodefault.golden"},
		{choices: []string{"foo", "bar", "baz"}, defValue: "bar", file: "testdata/choice-default.golden"},
		{choices: []string{"foo", "bar", "baz"}, defValue: "bak", file: "testdata/choice-default-new.golden"},
	}

	for i, test := range tests {
		app := libcmd.NewApp("app", "some brief description")
		app.Choice(test.choices, "choice", 'c', test.defValue, "One of %s.")

		if err := compareHelpOutput(app, []string{"-h"}, test.file); err != nil {
			t.Errorf("Case %d, %v", i, err)
		}
	}
}

func TestHelpChoiceTemplate(t *testing.T) {
	tests := []struct {
		template string
		file     string
	}{
		{template: "Values: %s.", file: "testdata/choice-template-normal.golden"},
		{template: "Without template.", file: "testdata/choice-template-none.golden"},
		{template: "", file: "testdata/choice-template-empty.golden"},
	}

	for i, test := range tests {
		app := libcmd.NewApp("app", "some brief description")
		app.Choice([]string{"a", "b", "c"}, "choice", 'c', "", test.template)

		if err := compareHelpOutput(app, []string{"-h"}, test.file); err != nil {
			t.Errorf("Case %d, %v", i, err)
		}
	}
}
