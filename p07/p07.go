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

var gVerbose bool

// Double-edged graph.
type Bag struct {
  name string
  contains map[string]int
  containedBy map[string]int
}

type RuleGraph struct {
  bags map[string]*Bag
}

func NewRuleGraph() *RuleGraph {
  g := new(RuleGraph)
  g.bags = make(map[string]*Bag, 768)
  return g
}

func NewBag(name string) *Bag {
  bag := new(Bag)
  bag.name = name
  bag.contains = make(map[string]int, 64)
  bag.containedBy = make(map[string]int, 64)
  return bag
}

func (g *RuleGraph) AddBag(name string) *Bag {
  bag := g.bags[name]
  if bag == nil {
    bag = NewBag(name)
    g.bags[name] = bag
  }
  return bag
}

// Add the rule "container contains N bags".
func (g *RuleGraph) AddRule(contained string, container string, N int) error {
  // Create an edge from contained to its new parent.
  containedBag := g.AddBag(contained)
  containedBag.containedBy[container] = N
  if containedBag.contains[container] != 0 {
    return errors.New(fmt.Sprintf(
      "RuleGraph.AddRule: direct cycle '%s' <-> '%s'", contained, container))
  }
  containerBag := g.AddBag(container)
  containerBag.contains[contained] = N
  if containerBag.containedBy[contained] != 0 {
    return errors.New(fmt.Sprintf(
      "RuleGraph.AddRule: direct cycle '%s' <-> '%s'", container, contained))
  }
  return nil
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
      if err = g.AddRule(containedColor, containerColor, num); err != nil {
        return err
      }
    }
  }
  return nil
}

// Traverse the contained-by closure beginning at the requested bag.
//
// Every bag visited is by extension capable of containing the input bag.
//
// The functions are called like f(contained, container, N) corresponding
// to rules "container contains N bags of contained" reachable from the first
// bag.
//
// fpre and fpost are called for pre-order and post-order traversal.
// If either is nil, it will not be called. If fpre returns false, traversal
// will not recurse. This may be used to detect and break cycles.
func (g *RuleGraph) TraverseContainedBy(bagName string,
                                        fpre func(string, string, int) bool,
                                        fpost func(string, string, int)) {
  bag := g.bags[bagName]
  if bag != nil {
    for container, num := range bag.containedBy {
      if fpre == nil || fpre(bag.name, container, num) {
        g.TraverseContainedBy(container, fpre, fpost)
      }
    }
    if fpost != nil {
      for container, num := range bag.containedBy {
        fpost(bag.name, container, num)
      }
    }
  }
}

// Traverse the contains close beginning at the requested bag.
//
// Every bag visited is by extension contained by the input bag.
//
// The functions are called like f(container, contained, N) corresponding
// to rules "container contains N bags of contained" reachable from the first
// bag.
//
// fpre and fpost are called for pre-order and post-order traversal.
// If either is nil, it will not be called. If fpre returns false, traversal
// will not recurse. This may be used to detect and break cycles.
func (g *RuleGraph) TraverseContains(bagName string,
                                     fpre func(string, string, int) bool,
                                     fpost func(string, string, int)) {
  bag := g.bags[bagName]
  if bag != nil {
    for contained, num := range bag.contains {
      if fpre == nil || fpre(bag.name, contained, num) {
        if gVerbose {
          fmt.Printf("RECURSING: %s => contains %d %s bags\n",
            bag.name, num, contained)
        }
        g.TraverseContains(contained, fpre, fpost)
      }
    }
  }
}

func (g *RuleGraph) visit_contains_dagsort(bag *Bag,
                                          stack []*Bag,
                                          visited map[string]bool) (
                                            []*Bag, map[string]bool) {
  if visited[bag.name] {
    return stack, visited
  }
  visited[bag.name] = true
  for name, _ := range bag.contains {
    stack, visited = g.visit_contains_dagsort(g.bags[name], stack, visited)
  }
  stack = append(stack, bag)
  return stack, visited
}

// Sort in contains order.
//
// That is, the first bags in the returned slice refer to bags which
// contain no other bags, and the last few bags contain the most bags.
// More formally, if index[B1] < index[B2], then B1 does not contain B2.
func (g *RuleGraph) SortContains() []*Bag {
  stack := make([]*Bag, 0, len(g.bags))
  visited := make(map[string]bool, len(g.bags))

  for _, bag := range g.bags {
    stack, visited = g.visit_contains_dagsort(bag, stack, visited)
  }

  return stack
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
  gVerbose = verbose
  graph, err := ReadRulesFromFile(input_path)
  if err != nil {
    return err
  }

  // Count the number of bags which may indirectly contain 'shiny gold' bags.
  mayContainGold := make(map[string]int, len(graph.bags))
  graph.TraverseContainedBy("shiny gold",
    func(contained string, container string, num int) bool{
      isNew := mayContainGold[container] == 0
      mayContainGold[container] = num
      if (verbose) {
        fmt.Printf("%s bags contain %d %s bags\n", container, num, contained)
      }
      return isNew
    }, nil)

  if (verbose) {
    fmt.Println("==============================================")
  }
  fmt.Printf("Bags which may eventually contain shiny gold bags: %d.\n",
    len(mayContainGold))
  if (verbose) {
    fmt.Println("==============================================")
  }

  // Count the number of bags which every bag must contain.
  containsDag := graph.SortContains()
  containsCount := make(map[string]int, len(graph.bags))
  for _, bag := range containsDag {
    for containedName, num := range bag.contains {
      containsCount[bag.name] += num * (1 + containsCount[containedName])
      if (verbose) {
        fmt.Printf("%s bags contain %d %s bags, each counting for %d, now %d\n",
          bag.name, num, containedName, containsCount[containedName],
          containsCount[bag.name])
      }
    }
    if verbose && len(bag.contains) == 0 {
      fmt.Printf("%s bags contain no other bags\n", bag.name)
    }
  }

  if (verbose) {
    fmt.Println("==============================================")
  }
  fmt.Printf("Number of bags which shiny gold must contain: %d.\n",
    containsCount["shiny gold"])
  if (verbose) {
    fmt.Println("==============================================")
  }

  return nil
}
