package p07

import (
  "io"
  "bufio"
  "os"
  "fmt"
  "strings"
  "errors"
  "strconv"
)

type ContainedBy struct {
  num int
  whom string
}

type RuleGraph struct {
  bags map[string][]ContainedBy
}

func NewRuleGraph() *RuleGraph {
  g := new(RuleGraph)
  g.bags = make(map[string][]ContainedBy, 768)
  return g
}

func (g *RuleGraph) AddBag(name string) {
  edges := g.bags[name]
  if edges == nil {
    edges = make([]ContainedBy, 0, 64)
    g.bags[name] = edges
  }
}

// Add the rule "container contains N bags".
func (g *RuleGraph) AddRule(bag string, container string, N int) {
  // Create an edge from bag to its new parent.
  g.AddBag(bag)
  g.AddBag(container)
  g.bags[bag] = append(g.bags[bag], ContainedBy{N, container})
}

// Parse a rule of the form:
// <color1> bags contain {<N> <color> bags[, ...]|no other bags}.
func (g *RuleGraph) ParseRule(rule string) error {
  parts := strings.SplitN(rule, " bags contain ", 2)
  containerColor := parts[0]
  rules := strings.Split(parts[1], ", ")
  if len(rules) == 0 || (len(rules) == 1 && strings.HasPrefix(rules[0], "no")) {
    g.AddBag(containerColor)
  } else {
    for _, rule := range rules {
      ruleParts := strings.SplitN(rule, " ", 2)
      num, err := strconv.Atoi(ruleParts[0])
      if err != nil {
        return err
      }
      contained := ruleParts[1]
      end := strings.Index(contained, " bag")
      if end < 0 {
        return errors.New(fmt.Sprintf("malformed bag list '%s'", contained))
      }
      containedColor := contained[:end]
      g.AddRule(containedColor, containerColor, num)
    }
  }
  return nil
}

// Number of bags which can directly contain the given bag.
func (g *RuleGraph) Order(bag string) int {
  if g.bags[bag] == nil {
    return -1
  }
  return len(g.bags[bag])
}

// Traverse the ContainedBy closure beginning at the requested bag.
//
// The function will be called like f(C, container, N) corresponding
// to rules "container contains N bags of color C" reachable from the first bag.
// (Therefore, every bag visited is by extension capable of containing the
// input bag.)
//
// If the function returns false, traversal will not recurse.
// Therefore, unless the callback somehow knows the graph does not contain
// cycles, it should track every node visited and prevent infinite loops.
func (g *RuleGraph) TraverseContained(bag string,
                                      f func (string, string, int) bool) {
  edges := g.bags[bag]
  if edges != nil {
    for _, edge := range edges {
      if f(bag, edge.whom, edge.num) {
        g.TraverseContained(edge.whom, f)
      }
    }
  }
}

func ReadRules(input io.Reader) (*RuleGraph, error) {
  scanner := bufio.NewScanner(input)
  scanner.Split(bufio.ScanLines)
  graph := NewRuleGraph()
  for scanner.Scan() {
    if err := graph.ParseRule(scanner.Text()); err != nil {
      return graph, err
    }
  }
  return graph, nil
}

func ReadRulesFromFile(path string) (*RuleGraph, error) {
  file, err := os.Open(path)
  if err != nil {
    return nil, err
  }
  defer file.Close()
  return ReadRules(file)
}

func Main(input_path string, verbose bool, args []string) error {
  graph, err := ReadRulesFromFile(input_path)
  if err != nil {
    return err
  }

  // Count the number of bags which may indirectly contain 'shiny gold' bags.
  mayContainGold := make(map[string]int, len(graph.bags))
  graph.TraverseContained("shiny gold",
    func(contained string, container string, num int) bool{
      isNew := mayContainGold[container] == 0
      mayContainGold[container] = num
      if (verbose) {
        fmt.Printf("%s bags contain %d %s bags\n", container, num, contained)
      }
      return isNew
    })

  fmt.Printf("Bags which may eventually contain shiny gold bags: %d.\n",
    len(mayContainGold))

  return nil
}
