package p02

import (
  "io"
  "fmt"
  "bufio"
  "os"
  "strings"
  "strconv"
  "errors"
)

type Password struct {
  x int
  y int
  char byte
  password string
}

// Part 1 validation
func (p *Password) P1Valid() bool {
  count := 0
  char := p.char
  for cidx := 0; cidx < len(p.password); cidx++ {
    if char == p.password[cidx] {
      count += 1
      if count > p.y {
        return false
      }
    }
  }
  return count >= p.x
}

// Part 2 validation
func (p *Password) P2Valid() bool {
  p1 := p.password[p.x - 1] == p.char
  p2 := p.password[p.y - 1] == p.char
  return p1 != p2
}

func parse_spec(spec string) (int, int, error) {
  parts := strings.Split(spec, "-")
  if len(parts) != 2 {
    return 0, 0, errors.New(fmt.Sprintf("invalid 'X-Y' specifier '%s'", spec))
  }
  var (x int
       y int
       err error)
  x, err = strconv.Atoi(parts[0])
  if err != nil {
    return 0, 0, err
  }
  y, err = strconv.Atoi(parts[1])
  if err != nil {
    return x, 0, err
  }
  return x, y, nil
}

func ParsePasswords(r io.Reader) ([]Password, error) {
  scanner := bufio.NewScanner(r)
  scanner.Split(bufio.ScanWords)

  passwords := make([]Password, 1000)
  count := 0
  for scanner.Err() == nil && scanner.Scan() {
    // X-Y char: password (three words)
    spec := scanner.Text()

    var char string
    if scanner.Scan() && scanner.Err() == nil {
      char = scanner.Text()
    } else {
      break
    }

    var password string
    if scanner.Scan() && scanner.Err() == nil {
      password = scanner.Text()
    } else {
      break
    }

    x, y, err := parse_spec(spec)
    if err != nil {
      return passwords, err
    }

    passwords[count] = Password{x, y, char[0], password}
    count++
  }

  return passwords, scanner.Err()
}

func ParsePasswordsFromFile(path string) ([]Password, error) {
  file, err := os.Open(path)
  if err != nil {
    return []Password{}, err
  }
  defer file.Close()
  return ParsePasswords(file)
}

type password_validator func(*Password) (bool)

func report_valid(passwords []Password, validate password_validator) int {
  lpass := len(passwords)
  nvalid := 0
  for _, password := range passwords {
    if validate(&password) {
      nvalid += 1
    }
  }

  if nvalid < lpass {
    fmt.Printf("Oh no! Only %d / %d passwords are valid!\n", nvalid, lpass)
  } else {
    fmt.Printf("Yay! All %d passwords are valid!\n", lpass)
  }
  return nvalid
}

func Main(input_path string, verbose bool, args []string) error {
  passwords, err := ParsePasswordsFromFile(input_path)
  if err != nil {
    return err
  }

  fmt.Printf("Validation for part 1:\n  ")
  _ = report_valid(passwords, (*Password).P1Valid)
  fmt.Printf("Validation for part 2:\n  ")
  _ = report_valid(passwords, (*Password).P2Valid)

  return nil
}
