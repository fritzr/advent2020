package p19

import (
  "fmt"
  "errors"
  "strings"
  "strconv"
  "github.com/fritzr/advent2020/util"
)

var gVerbose = false

type Rule interface { }

// Match a literal string.
type Literal struct {
  literal string
}

// Match a sequence of rules in order.
type Sequence struct {
  all []int
}

// Match any of a set of rules.
type Selector struct {
  any []Rule
}

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


// This function returns the prefix of text which is matched by repeating
// the given sequence of rules.
//func (g *Grammar) ruleAcceptsMany(rules Rule, text string) (bool, string) {
//}

func (g *Grammar) literal(id int, rule *Literal, text string, index int) []int {
  // Match literal prefix.
  literal := rule.literal
  if index < len(text) && text[index:index + len(literal)] == literal {
    return []int{index + len(literal)}
  }
  return []int{}
}

func (g *Grammar) any(id int, rule *Selector, text string, index int) []int {
  // Match any of a set of other rules.
  valid := make([]int, 0, len(rule.any))
  for _, subRule := range rule.any {
    valid = append(valid, g.all(id, subRule.(*Sequence), text, index)...)
  }
  return valid
}

func (g *Grammar) all(id int, rule *Sequence, text string, index int) []int {
  // Match all sub-rules sequentially.
  var suffixes = []int{index}
  for _, ruleIndex := range rule.all {
    suffixes = g.filter(ruleIndex, g.rules[ruleIndex], text, suffixes, g.prefix)
    // We can stop trying if there are no valid suffixes anymore.
    if len(suffixes) == 0 {
      break
    }
  }
  return suffixes
}

func (g *Grammar) filter(id int, rule Rule, text string, suffixes []int,
    accept func(id int, rule Rule, text string, index int) []int) []int {
  valid := make([]int, 0)
  for _, index := range suffixes {
    valid = append(valid, accept(id, rule, text, index)...)
  }
  return valid
}

func (g *Grammar) prefix(id int, rule Rule, text string, index int) []int {
  switch r := rule.(type) {
    case *Literal: return g.literal(id, r, text, index)
    case *Selector: return g.any(id, r, text, index)
    case *Sequence: return g.all(id, r, text, index)
    default: panic("unhandled rule type")
  }
}

func (g *Grammar) Accepts(text string) bool {
  if gVerbose {
    fmt.Printf("**** (len=%d) %s\n", len(text), text)
  }
  suffixes := g.prefix(0, g.rules[0], text, 0)
  if len(suffixes) > 0 {
    for _, suffix := range suffixes {
      if gVerbose {
        fmt.Printf("    text[:%d] matches (%s).\n", suffix, text[:suffix])
      }
      if suffix >= len(text) {
        if gVerbose {
          fmt.Println("=> Accepted.")
        }
        return true
      }
    }
  } else if gVerbose {
    fmt.Println("    no matches.")
  }
  if gVerbose {
    fmt.Println("=> Rejected.")
  }
  return false
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
  gVerbose = verbose
  groups, err := util.ReadLineGroupsFromFile(input_path)
  if err != nil {
    return err
  }

  if len(groups) != 2 {
    return errors.New("expected 2 groups in input")
  }

  g := NewGrammar()
  err = g.ParseRules(groups[0])
  if err != nil {
    return err
  }

  // Part 1: see how many messages are accepted.
  messages := strings.Split(groups[1], "\n")
  valid := 0
  for _, message := range messages {
    if g.Accepts(message) {
      valid++
    }
  }
  fmt.Printf("%d / %d messages are valid.\n", valid, len(messages))

  // Part 2: replace 8 and 11 with some recursive rules.
  valid = 0
  g.SetRule(8, &Selector{[]Rule{
    &Sequence{[]int{42}}, &Sequence{[]int{42, 8}}}})
  g.SetRule(11, &Selector{[]Rule{
    &Sequence{[]int{42, 31}}, &Sequence{[]int{42, 11, 31}}}})
  for _, message := range messages {
    if g.Accepts(message) {
      valid++
    }
  }
  fmt.Printf("%d / %d messages are valid with recursive rules.\n",
    valid, len(messages))

  return nil
}
