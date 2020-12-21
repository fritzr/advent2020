package p19

import (
  "io"
  "fmt"
  "errors"
  "strings"
  "strconv"
  "github.com/fritzr/advent2020/util"
)

const RULE_LITERAL = 1
const RULE_ANY = 2
const RULE_ALL = 3

// One of the RULE_* constants.
type Rule interface {
  Type() int
}

// Match a literal string.
type Literal struct {
  literal string
}

func (r *Literal) Type() int {
  return RULE_LITERAL
}

// Match a sequence of rules in order.
type Sequence struct {
  all []int
}

func (r *Sequence) Type() int {
  return RULE_ALL
}

// Match any of a set of rules.
type Selector struct {
  any []Rule
}

func (r *Selector) Type() int {
  return RULE_ANY
}

// It might be better to write a state machine, but RDP is easier.
type Grammar struct {
  rules map[int]Rule
}

func NewGrammar() *Grammar {
  g := new(Grammar)
  g.rules = make(map[int]Rule)
  return g
}

// We support two specific types of recursive rules:
//   1. B -> A | AB
//      This means "B matches A 1 or more times".
//   2. B -> AC | ABC
//      This means "B matches 1 or more A followed by the same number of C".


// For 1: match 1 or more copies of the sequence.
func (g *Grammar) many(id int, sequence *Sequence, text string) (bool, string) {
  // Consume the sequence 1 or more times.
  valid, suffix := g.all(id, sequence, text)
  validOnce := valid
  lastValidSuffix := suffix
  for valid {
    lastValidSuffix = suffix
    valid, suffix = g.all(id, sequence, suffix)
  }
  return validOnce, lastValidSuffix
}

// Match exactly N copies of the sequence.
func (g *Grammar) matchN(id int, sequence *Sequence, text string, n int) (
    bool, string) {
  valid := true
  suffix := text
  for ; valid && n > 0; n-- {
    valid, suffix = g.all(id, sequence, suffix)
  }
  return valid, suffix
}

// For 2:
// This function returns the prefix of text which is matched by matching
// some number of copies of lhs followed by the same number of copies of rhs.
func (g *Grammar) symmetric(id int, lhs *Sequence, rhs *Sequence, text string) (
    bool, string) {
  // Incrementally match symmetric patterns of (n*lhs) + (n*rhs) until we get
  // a match. TODO do we need to get *every* such match?
  for n := 1; n < len(text) / 2; n++ {
    lValid, lSuffix := g.matchN(id, lhs, text, n)
    if lValid {
      rValid, rSuffix := g.matchN(id, rhs, lSuffix, n)
      if rValid {
        return true, rSuffix
      }
    }
  }
  return false, text
}

func (g *Grammar) literal(id int, rule *Literal, text string) (bool, string) {
  // Match literal prefix.
  literal := rule.literal
  if strings.HasPrefix(text, literal) {
    return true, text[len(literal):]
  }
  return false, text
}

func (g *Grammar) any(id int, rule *Selector, text string) (bool, string) {
  // Match any of a set of other rules.

  // First, check if the rule forms a recursive relation that we support.
  recursionPoint := -1
  lastRule := rule.any[len(rule.any)-1].(*Sequence).all
  for subIndex, subid := range lastRule {
    if subid == id {
      if recursionPoint != -1 {
        panic("multiple recursion unsupported")
      }
      recursionPoint = subIndex
      // break // commented out so we will check the assertion above
    }
  }

  if recursionPoint == len(lastRule) - 1 {
    // 1. matches many: B -> A | AB
    return g.many(id, &Sequence{lastRule[:recursionPoint]}, text)
  } else if recursionPoint > 1 /* && XXX? */{
    // 2. matches symmetric: B -> AC | ABC
    return g.symmetric(id, &Sequence{lastRule[:recursionPoint]},
      &Sequence{lastRule[recursionPoint + 1:]}, text)
  } else if recursionPoint >= 0 {
    // This isn't to say we will catch all unsupported recursive relations;
    // specifically, we will not catch indirect recursion!
    panic("unsupported recursive relation")
  }

  // No recursion. Simply try each sub-rule.
  for _, subRule := range rule.any {
    valid, suffix := g.prefix(id, subRule, text)
    if valid {
      return true, suffix
    }
  }
  return false, text
}

func (g *Grammar) all(id int, rule *Sequence, text string) (bool, string) {
  // Match all sub-rules sequentially.
  suffix := text
  valid := false
  for _, ruleIndex := range rule.all {
    valid, suffix = g.prefix(ruleIndex, g.rules[ruleIndex], suffix)
    if !valid {
      return false, suffix
    }
  }
  return true, suffix
}

func (g *Grammar) prefix(id int, rule Rule, text string) (bool, string) {
  switch rule.Type() {
    case RULE_LITERAL: return g.literal(id, rule.(*Literal), text)
    case RULE_ANY: return g.any(id, rule.(*Selector), text)
    case RULE_ALL: return g.all(id, rule.(*Sequence), text)
    default: panic("unhandled rule type")
  }
}

func (g *Grammar) Accepts(text string) bool {
  valid, tail := g.prefix(0, g.rules[0], text)
  return valid && tail == ""
}

func (g *Grammar) SetRule(ruleId int, rule Rule) {
  g.rules[ruleId] = rule
}

func (g *Grammar) ParseRule(rule string) error {
  parts := strings.Split(rule, ": ")
  if len(parts) != 2 {
    return errors.New("expected 'index: rule'")
  }

  index, err := strconv.Atoi(parts[0])
  if err != nil {
    return err
  }

  fields := strings.Fields(parts[1])
  if len(fields) == 0 {
    return errors.New("empty rule body")
  }

  // Literal rule.
  firstWord := fields[0]
  if len(fields) == 1 && len(firstWord) > 2 && (
      firstWord[0] == '"' && firstWord[len(firstWord)-1] == '"') {
    g.SetRule(index, &Literal{firstWord[1:len(firstWord)-1]})
    return nil
  }

  // Selector (A B | C D | ...) or Sequence (A B C D...).
  startIndex := 0
  rules := make([]Rule, 0)
  for fieldIndex := 0; fieldIndex < len(fields); fieldIndex++ {
    if fields[fieldIndex] == "|" {
      if fieldIndex == 0 || startIndex == fieldIndex {
        return errors.New(fmt.Sprintf(
          "rule parse error: unexpected '|' in '%s'", parts[1]))
      }
      elements, err2 := util.FieldsToInts(fields[startIndex:fieldIndex])
      if err2 != nil {
        return err2
      }
      rules = append(rules, &Sequence{all: elements})
      startIndex = fieldIndex + 1
    }
  }
  finalElements, err2 := util.FieldsToInts(fields[startIndex:])
  if err2 != nil {
    return err2
  }

  if len(rules) == 0 {
    g.SetRule(index, &Sequence{all: finalElements})
  } else {
    rules = append(rules, &Sequence{all: finalElements})
    g.SetRule(index, &Selector{any: rules})
  }
  return nil
}

func (g *Grammar) ParseRules(rulesText string) error {
  lines := strings.Split(rulesText, "\n")
  for _, line := range lines {
    if err := g.ParseRule(line); err != nil {
      return err
    }
  }
  return nil
}

func Main(input_path string, verbose bool, args []string) error {
  groups, err := util.ReadFile(input_path,
    func(input io.Reader) (interface{}, error) {
      return util.ScanInput(input, util.ScanLineGroups)
    })
  if err != nil {
    return err
  }

  textGroups := groups.([]string)
  if len(textGroups) != 2 {
    return errors.New("expected 2 groups in input")
  }

  g := NewGrammar()
  err = g.ParseRules(textGroups[0])
  if err != nil {
    return err
  }

  messages := strings.Split(textGroups[1], "\n")
  valid := 0
  for _, message := range messages {
    if g.Accepts(message) {
      valid++
    }
  }

  fmt.Printf("%d / %d messages are valid.\n", valid, len(messages))
  return nil
}
