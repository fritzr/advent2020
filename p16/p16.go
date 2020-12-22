package p16

import (
  "io"
  "fmt"
  "errors"
  "strconv"
  "strings"
  "github.com/fritzr/advent2020/util"
)

var gVerbose = false

type TicketField struct {
  name string
  ranges [][2]int
}

// Return the index of the first valid field in the ticket.
func (f *TicketField) FindValidField(ticket []int) int {
  for index, field := range ticket {
    if f.IsValid(field) {
      return index
    }
  }
  return -1
}

// Whether a field is valid.
func (f *TicketField) IsValid(value int) bool {
  for _, extents := range f.ranges {
    if extents[0] <= value && value <= extents[1] {
      return true
    }
  }
  return false
}

func parseTicketField(line string) (*TicketField, error) {
  // name: x-y or a-b...
  lr := strings.Split(line, ": ")
  if len(lr) != 2 {
    return nil, errors.New("expected field name separator")
  }
  name := lr[0]
  rangeStrings := strings.Split(lr[1], " or ")
  if len(rangeStrings) == 0 {
    return nil, errors.New("expected non-empty field range")
  }
  ranges := make([][2]int, len(rangeStrings))
  for index, rangeString := range rangeStrings {
    rangeFields := strings.Split(rangeString, "-")
    if len(rangeFields) != 2 {
      return nil, errors.New("expected range separator")
    }
    min, err := strconv.Atoi(rangeFields[0])
    if err != nil {
      return nil, err
    }
    max, err := strconv.Atoi(rangeFields[1])
    if err != nil {
      return nil, err
    }
    ranges[index] = [2]int{min, max}
  }
  return &TicketField{name, ranges}, nil
}

func parseTicketFields(fieldsText string) (fields []*TicketField, err error) {
  // One field per line.
  lines := strings.Split(fieldsText, "\n")
  fields = make([]*TicketField, len(lines))
  for lineNumber, line := range lines {
    fields[lineNumber], err = parseTicketField(line)
    if err != nil {
      break
    }
  }
  return fields, err
}

func ReadListOfInts(text string, sep string) (list [][]int, err error) {
  // Each line is a list of integers.
  lines := strings.Split(text, "\n")
  list = make([][]int, len(lines))
  for lineNumber, line := range lines {
    list[lineNumber], err = util.FieldsToInts(strings.Split(line, sep))
    if err != nil {
      break
    }
  }
  return list, err
}

func parseTickets(ticketText string) ([][]int, error) {
  // Skip the first line.
  eolIndex := strings.IndexByte(ticketText, '\n')
  if eolIndex < 0 {
    return nil, errors.New("expected end-of-line in ticket")
  }
  // Obtain list of integer lists from the subsequent lines.
  return ReadListOfInts(ticketText[eolIndex+1:], ",")
}

func findValidTickets(fields []*TicketField, tickets [][]int) (
    validTicketNumbers util.Set, errorRate int) {
  validTicketNumbers = make(util.Set)
  for ticketNumber, ticket := range tickets {
    if gVerbose {
      fmt.Printf("[%d] ## %s\n", ticketNumber, util.ArrayString(ticket))
    }
    allValid := true
    for _, value := range ticket {
      valid := false
      for _, fieldSpec := range fields {
        if fieldSpec.IsValid(value) {
          if gVerbose {
            fmt.Printf("     %d is valid for %s\n", value, fieldSpec.name)
          }
          valid = true
          break
        } else if gVerbose {
          fmt.Printf("     %d is invalid for %s\n", value, fieldSpec.name)
        }
      }
      if !valid {
        allValid = false
        if gVerbose {
          fmt.Printf("[%d] is invalid\n", ticketNumber)
        }
        errorRate += value
      }
    }
    if allValid {
      if gVerbose {
        fmt.Printf("[%d] is valid\n", ticketNumber)
      }
      validTicketNumbers[ticketNumber] = true
    }
  }
  return validTicketNumbers, errorRate
}

func identifyFields(fields []*TicketField, tickets [][]int, ticketIndexes util.Set)(
    []string, error) {
  // Now identify which fields are which based on validity.
  // This structure maps field names to the field indexes which are possible.
  // Whenever a range constraint for that field is violated, we remove it from
  // the set of possible indexes.
  possibleFieldIndexes := make(map[string]util.Set)
  for _, field := range fields {
    possibleFieldIndexes[field.name] = make(util.Set)
    for index := 0; index < len(fields); index++ {
      possibleFieldIndexes[field.name][index] = true
    }
  }

  // Mask out the fields which are invalid by index.
  for ticketNumber, _ := range ticketIndexes {
    ticket := tickets[ticketNumber]
    for fieldIndex, field := range ticket {
      for _, fieldSpec := range fields {
        if !fieldSpec.IsValid(field) {
          delete(possibleFieldIndexes[fieldSpec.name], fieldIndex)
          if gVerbose {
            fmt.Printf(
              "[%d]: %s cannot be [%d] because %d is invalid (now: %v)\n",
              ticketNumber, fieldSpec.name, fieldIndex, field,
              possibleFieldIndexes[fieldSpec.name])
          }
        }
      }
    }
  }

  // Now reduce to a one-to-one map of field names to indexes.
  // Greedily fix any values which have one and only one possibility.
  fieldNames := make([]string, len(fields))
  previouslySelected := -1
  selectedFields := 0
  nFields := len(fieldNames)
  for selectedFields != previouslySelected && selectedFields != nFields {
    previouslySelected = selectedFields
    // Select fields which have unique choices.
    selectedNow := make(map[int]string)
    for fieldName, possibleIndexes := range possibleFieldIndexes {
      // Select the field if it is unique.
      if len(possibleIndexes) == 1 {
        for index, _ := range(possibleIndexes) {
          if gVerbose {
            fmt.Printf("selecting field [%d] for %s\n", index, fieldName)
          }
          // Check we didn't just select it for a different field.
          if selectedNow[index] != "" {
            return fieldNames, errors.New(fmt.Sprintf(
              "'%s' and '%s' both uniquely match field [%d]",
              fieldName, selectedNow[index], index))
          }
          selectedNow[index] = fieldName
        }
      }
    }
    // Remove the selected fields as possibilities for anyone else.
    for selectedIndex, fieldName := range selectedNow {
      for oldName, possibleIndexes := range possibleFieldIndexes {
        if gVerbose && possibleIndexes[selectedIndex] {
          fmt.Printf("[%d] was assigned to %s: no longer a candidate for %s\n",
            selectedIndex, fieldName, oldName)
        }
        delete(possibleIndexes, selectedIndex)
      }
      // Then commit the selected index.
      fieldNames[selectedIndex] = fieldName
    }
    selectedFields += len(selectedNow)
  }

  // All possibleFields left are errors!
  if selectedFields < nFields {
    var errorStr strings.Builder
    for fieldName, indexes := range possibleFieldIndexes {
      if len(indexes) == 0 {
        errorStr.WriteString(fmt.Sprintf(
          "no matches for field '%s'\n", fieldName))
      } else {
        errorStr.WriteString(fmt.Sprintf(
          "multiple (%d) matches for field '%s': %v\n",
          len(indexes), fieldName, indexes))
      }
    }
    return fieldNames, errors.New(errorStr.String())
  }

  return fieldNames, nil
}


func Main(input_path string, verbose bool, args []string) error {
  gVerbose = verbose
  groups, err := util.ReadFile(input_path,
    func(input io.Reader) (interface{}, error) {
      return util.ScanInput(input, util.ScanLineGroups)
    })
  if err != nil {
    return err
  }

  strGroups := groups.([]string)
  if len(strGroups) != 3 {
    return errors.New("invalid input format")
  }

  // Input is broken into {fields, my ticket, nearby tickets}.
  fields, err := parseTicketFields(strGroups[0])
  if err != nil {
    return err
  }

  if verbose {
    for _, field := range fields {
      fmt.Printf("%s:", field.name)
      for _, extents := range field.ranges {
        fmt.Printf("  %d-%d", extents[0], extents[1])
      }
      fmt.Println()
    }
  }

  otherTickets, err := parseTickets(strGroups[2])
  if err != nil {
    return err
  }

  myTickets, err := parseTickets(strGroups[1])
  if err != nil {
    return err
  }

  // Part 1: filter out invalid tickets.
  validTickets, errorRate := findValidTickets(fields, otherTickets)

  fmt.Printf("Ticket scanning error rate: %d\n", errorRate)
  fmt.Printf("There were %d valid tickets.\n", len(validTickets))

  // Part 2: find ticket fields. Assume our ticket(s) are valid.

  if verbose {
    fmt.Printf("My tickets: ")
    for _, myTicket := range myTickets {
      fmt.Printf("# %s\n", util.ArrayString(myTicket))
    }
  }

  allTickets := append(otherTickets, myTickets...)
  validTickets[len(allTickets)-1] = true
  fieldNames, err := identifyFields(fields, allTickets, validTickets)
  if err != nil {
    return err
  }

  // Identify the fields on my ticket(s).
  fmt.Println("My ticket fields:")
  for _, myTicket := range myTickets {
    departureProd := 1
    for index, name := range fieldNames {
      myValue := myTicket[index]
      fmt.Printf("  [%2d] %s: %d\n", index, name, myValue)
      if strings.HasPrefix(name, "departure") {
        departureProd *= myValue
      }
    }
    fmt.Printf("Product of 'departure' fields: %d\n", departureProd)
  }

  return nil
}
