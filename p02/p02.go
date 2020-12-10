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
  min int
  max int
  char byte
  password string
}

func (p *Password) Valid() bool {
  count := 0
  char := p.char
  for cidx := 0; cidx < len(p.password); cidx++ {
    if char == p.password[cidx] {
      count += 1
      if count > p.max {
        return false
      }
    }
  }
  return count >= p.min
}

func parse_spec(spec string) (int, int, error) {
  parts := strings.Split(spec, "-")
  if len(parts) != 2 {
    return 0, 0, errors.New(fmt.Sprintf("invalid 'min-max' specifier '%s'", spec))
  }
  var (min int
       max int
       err error)
  min, err = strconv.Atoi(parts[0])
  if err != nil {
    return 0, 0, err
  }
  max, err = strconv.Atoi(parts[1])
  if err != nil {
    return min, 0, err
  }
  return min, max, nil
}

func ParsePasswords(r io.Reader) ([]Password, error) {
  scanner := bufio.NewScanner(r)
  scanner.Split(bufio.ScanWords)

  passwords := make([]Password, 1000)
  count := 0
  for scanner.Err() == nil && scanner.Scan() {
    // min-max char: password (three words)
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

    min, max, err := parse_spec(spec)
    if err != nil {
      return passwords, err
    }

    passwords[count] = Password{min, max, char[0], password}
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

func Main(input_path string, verbose bool, args []string) error {
  passwords, err := ParsePasswordsFromFile(input_path)
  if err != nil {
    return err
  }

  nvalid := 0
  for _, password := range passwords {
    if password.Valid() {
      nvalid += 1
    }
  }

  lpass := len(passwords)
  if nvalid < lpass {
    fmt.Printf("Oh no! Only %d / %d passwords are valid!\n", nvalid, lpass)
  } else {
    fmt.Printf("Yay! All %d passwords are valid!\n", lpass)
  }

  return nil
}
