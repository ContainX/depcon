package envsubst

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

var theWordIsGo = map[string]string{
	"WORD": "go",
}

var funWithDigits = map[string]string{
	"FUN1": "gopher",
}

// Wrap substituteVariableReferences for easy testing
func subst(input string, vars map[string]string) string {
	return substitute(input, false, vars)
}

func substitute(input string, preserveUndef bool, vars map[string]string) string {
	return Substitute(strings.NewReader(input), preserveUndef, func(s string) string {
		return vars[s]
	})
}

func TestEmptyInput(t *testing.T) {
	result := subst("", theWordIsGo)
	assert.Equal(t, result, "")
}

func TestNoVariables(t *testing.T) {
	result := subst("hello world", theWordIsGo)
	assert.Equal(t, "hello world", result)
}

func TestSimpleVariable(t *testing.T) {
	result := subst("hello $WORD world", theWordIsGo)
	assert.Equal(t, "hello $WORD world", result)

	result = subst("hello ${WORD} world", theWordIsGo)
	assert.Equal(t, "hello go world", result)
}

func TestSimpleVariableAtStart(t *testing.T) {
	result := subst("$WORD home world", theWordIsGo)
	assert.Equal(t, "$WORD home world", result)

	result = subst("${WORD} home world", theWordIsGo)
	assert.Equal(t, "go home world", result)
}

func TestSimpleVariableAtEnd(t *testing.T) {
	result := subst("let's $WORD", theWordIsGo)
	assert.Equal(t, "let's $WORD", result)

	result = subst("let's ${WORD}", theWordIsGo)
	assert.Equal(t, "let's go", result)
}

func TestOnlyVariable(t *testing.T) {
	result := subst("$WORD", theWordIsGo)
	assert.Equal(t, "$WORD", result)

	result = subst("${WORD}", theWordIsGo)
	assert.Equal(t, result, "go")
}

func TestRunOnVariable(t *testing.T) {
	result := subst("$WORD$WORD$WORD!", theWordIsGo)
	assert.Equal(t, "$WORD$WORD$WORD!", result)

	result = subst("${WORD}${WORD}${WORD}!", theWordIsGo)
	assert.Equal(t, "gogogo!", result)
}

func TestRunOnVariableWithNonVariableTextPrefix(t *testing.T) {
	result := subst("$WORD,no$WORD", theWordIsGo)
	assert.Equal(t, "$WORD,no$WORD", result)

	result = subst("${WORD},no${WORD}", theWordIsGo)
	assert.Equal(t, "go,nogo", result)
}

func TestSimpleStandAloneDollar(t *testing.T) {
	result := subst("2 $ for your $WORD thoughts", theWordIsGo)
	assert.Equal(t, "2 $ for your $WORD thoughts", result)

	result = subst("2 ${} for your $WORD thoughts", theWordIsGo)
	assert.Equal(t, "2 ${} for your $WORD thoughts", result)
}

func TestSimpleStandAloneDollarAtStart(t *testing.T) {
	result := subst("$ for your $WORD thoughts", theWordIsGo)
	assert.Equal(t, "$ for your $WORD thoughts", result)

	result = subst("${} for your $WORD thoughts", theWordIsGo)
	assert.Equal(t, "${} for your $WORD thoughts", result)
}

func TestSimpleStandAloneDollarAtEnd(t *testing.T) {
	result := subst("$WORD, find some $", theWordIsGo)
	assert.Equal(t, result, "$WORD, find some $")

	result = subst("${WORD}, find some ${}", theWordIsGo)
	assert.Equal(t, "go, find some ${}", result)
}

func TestOnlyStandAloneDollar(t *testing.T) {
	result := subst("$", theWordIsGo)
	assert.Equal(t, "$", result)

	result = subst("${}", theWordIsGo)
	assert.Equal(t, "${}", result)
}

func TestStandAloneDollarSuffix(t *testing.T) {
	result := subst("$WORD$", theWordIsGo)
	assert.Equal(t, "$WORD$", result)

	result = subst("${WORD}${}", theWordIsGo)
	assert.Equal(t, "go${}", result)
}

func TestBracingStartsMidtoken(t *testing.T) {
	result := subst("$WORD{FISH}", theWordIsGo)
	assert.Equal(t, "$WORD{FISH}", result)

	result = subst("$WORD{WORD}", theWordIsGo)
	assert.Equal(t, "$WORD{WORD}", result)

	result = subst("$WO{RD}", theWordIsGo)
	assert.Equal(t, "$WO{RD}", result)
}

func TestUnclosedBraceWithPrefix(t *testing.T) {
	result := subst("$WORD{WORD", theWordIsGo)
	assert.Equal(t, "$WORD{WORD", result)
}

func TestUnclosedBracedDollar(t *testing.T) {
	result := subst("${", theWordIsGo)
	assert.Equal(t, "${", result)
}

func TestUnclosedBracedDollarAndSubsequentTokens(t *testing.T) {
	result := subst("${ stuff }", theWordIsGo)
	assert.Equal(t, "${ stuff }", result)
}

func TestUnclosedBracedDollarWithSuffix(t *testing.T) {
	result := subst("${WORD", theWordIsGo)
	assert.Equal(t, "${WORD", result)
}

func TestUnclosedBracedDollarWithSuffixAndSubsequentTokens(t *testing.T) {
	result := subst("no ${WORD yet", theWordIsGo)
	assert.Equal(t, "no ${WORD yet", result)
}

func TestDoubleBracedDollar(t *testing.T) {
	result := subst("this ${{WORD}} should not be touched", theWordIsGo)
	assert.Equal(t, "this ${{WORD}} should not be touched", result)
}

func TestUndefinedVariablesAreRemovedByDefault(t *testing.T) {
	result := subst("nothing $WORD2 to see here", theWordIsGo)
	assert.Equal(t, "nothing $WORD2 to see here", result)

	result = subst("nothing ${WORD2} to see here", theWordIsGo)
	assert.Equal(t, "nothing  to see here", result)

	result = subst("nothing${WORD2}to see here", theWordIsGo)
	assert.Equal(t, "nothingto see here", result)
}

func TestUndefinedVariablesArePreservedWhenWanted(t *testing.T) {
	result := substitute("nothing $WORD2 to see here", true, theWordIsGo)
	assert.Equal(t, "nothing $WORD2 to see here", result)

	result = substitute("nothing ${WORD2} to see here", true, theWordIsGo)
	assert.Equal(t, "nothing ${WORD2} to see here", result)

	result = substitute("nothing${WORD2}to see here", true, theWordIsGo)
	assert.Equal(t, "nothing${WORD2}to see here", result)
}

func TestVariableNamesDontBeginWithADigit(t *testing.T) {
	result := subst("$1A", theWordIsGo)
	assert.Equal(t, "$1A", result)

	result = subst("${1A}", theWordIsGo)
	assert.Equal(t, "${1A}", result)
}

func TestVariableNamesAllowDigitsAfterFirstCharacter(t *testing.T) {
	result := subst("$FUN1", funWithDigits)
	assert.Equal(t, "$FUN1", result)

	result = subst("${FUN1}", funWithDigits)
	assert.Equal(t, "gopher", result)

	result = subst("no $FUN2", funWithDigits)
	assert.Equal(t, "no $FUN2", result)

	result = subst("no ${FUN2}", funWithDigits)
	assert.Equal(t, "no ", result)
}
