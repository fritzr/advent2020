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

type Rule interface {
  Type() int // one of the RULE_* constants
}

type Literal struct {
  literal string
}

func (r *Literal) Type() int {
  return RULE_LITERAL
}

type Sequence struct {
  all []int // index reference to other rules
}

func (r *Sequence) Type() int {
  return RULE_ALL
}

type Selector struct {
  any []Rule // collection of references to other rules
}

func (r *Selector) Type() int {
  return RULE_ANY
}

type Grammar struct {
  rules map[int]Rule
}

func NewGrammar() *Grammar {
  g := new(Grammar)
  g.rules = make(map[int]Rule)
  return g
}

func (g *Grammar) ruleAcceptsPrefix(rule Rule, text string) (bool, string) {
  switch rule.Type() {
    case RULE_LITERAL:
      // Match literal prefix.
      literal := rule.(*Literal).literal
      if strings.HasPrefix(text, literal) {
        return true, text[len(literal):]
      }
      return false, text

    case RULE_ANY:
      // Match any of a set of other rules.
      for _, r := range rule.(*Selector).any {
        valid, suffix := g.ruleAcceptsPrefix(r, text)
        if valid {
          return true, suffix
        }
      }
      return false, text

    case RULE_ALL:
      // Match all sub-rules sequentially.
      suffix := text
      valid := false
      for _, ruleIndex := range rule.(*Sequence).all {
        valid, suffix = g.ruleAcceptsPrefix(g.rules[ruleIndex], suffix)
        if !valid {
          return false, suffix
        }
      }
      return true, suffix

    default:
      panic("unhandled rule type")
  }
}

func (g *Grammar) Accepts(text string) bool {
  valid, tail := g.ruleAcceptsPrefix(g.rules[0], text)
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
