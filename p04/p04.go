package p04

import (
  "io"
  "bufio"
  "errors"
  "strings"
  "fmt"
  "strconv"
  "github.com/fritzr/advent2020/util"
)

type Validator func(value string) bool

func validate_byr(value string) bool {
  byr, err := strconv.Atoi(value)
  return len(value) == 4 && err == nil && byr >= 1920 && byr <= 2002
}

func validate_iyr(value string) bool {
  iyr, err := strconv.Atoi(value)
  return len(value) == 4 && err == nil && iyr >= 2010 && iyr <= 2020
}

func validate_eyr(value string) bool {
  eyr, err := strconv.Atoi(value)
  return len(value) == 4 && err == nil && eyr >= 2020 && eyr <= 2030
}

func validate_hgt(value string) bool {
  cm_idx := strings.Index(value, "cm")
  if cm_idx == 0 { return false }
  if cm_idx > 0 {
    if value[cm_idx:] != "cm" { return false }
    hgt_cm, err := strconv.Atoi(value[:cm_idx])
    return err == nil && hgt_cm >= 150 && hgt_cm <= 193
  }
  in_idx := strings.Index(value, "in")
  if in_idx <= 0 { return false }
  if value[in_idx:] != "in" { return false }
  in_cm, err := strconv.Atoi(value[:in_idx])
  return err == nil && in_cm >= 59 && in_cm <= 76
}

func validate_hcl(value string) bool {
  return len(value) == 7 && value[0] == '#' && (
    util.StringIsSubset(value[1:], "0123456789abcdef"))
}

func validate_ecl(value string) bool {
  return 1 == map[string]int{
    "amb":1, "blu":1, "brn":1, "gry":1, "grn":1, "hzl":1, "oth":1}[value]
}

func validate_pid(value string) bool {
  _, err := strconv.Atoi(value)
  return len(value) == 9 && err == nil
}

/*
func validate_cid(value string) bool {
  return true
}
*/

var validators = map[string]Validator{
  "byr": validate_byr,
  "iyr": validate_iyr,
  "eyr": validate_eyr,
  "hgt": validate_hgt,
  "hcl": validate_hcl,
  "ecl": validate_ecl,
  "pid": validate_pid /*, "cid": validate_cid */ }

type Passport struct {
  fields map[string]string
}

func (p *Passport) Present() bool {
  for key, _ := range validators {
    if p.fields[key] == "" {
      return false
    }
  }
  return true
}

func (p *Passport) Valid() bool {
  for key, validator := range validators {
    if p.fields[key] == "" || !validator(p.fields[key]) {
      return false
    }
  }
  return true
}

func split_word(word string) (string, string, error) {
  fields := strings.Split(word, ":")
  if len(fields) != 2 {
    return "", "", errors.New(fmt.Sprintf("invalid field '%s'", word))
  }
  return fields[0], fields[1], nil
}

func (p *Passport) Read(data string) error {
  // Split all the words in the data.
  words := bufio.NewScanner(strings.NewReader(data))
  words.Split(bufio.ScanWords)
  for words.Scan() {
    key, value, err := split_word(words.Text())
    if err != nil {
      return err
    }
    p.fields[key] = value
  }
  return nil
}

func NewPassport(data string) (Passport, error) {
  p := Passport{}
  p.fields = make(map[string]string, 8)
  return p, p.Read(data)
}

func ReadPassports(input io.Reader) ([]Passport, error) {
  scanner := bufio.NewScanner(input)
  scanner.Split(util.ScanLineGroups)

  passports := make([]Passport, 0, 1056)
  for scanner.Scan() {
    p, err := NewPassport(scanner.Text())
    if err != nil {
      return passports, err
    }
    passports = append(passports, p)
  }

  return passports, scanner.Err()
}

func Main(input_path string, verbose bool, args []string) error {
  iPassports, err := util.ReadFile(input_path,
    func(input io.Reader) (interface{}, error) {
      return ReadPassports(input)
    })
  if err != nil {
    return err
  }

  passports := iPassports.([]Passport)
  npresent := 0
  nvalid := 0
  for _, p := range passports {
    if p.Present() {
      npresent++
      if p.Valid() {
        nvalid++
      }
    }
  }

  fmt.Printf("%d / %d passports have all fields\n", npresent, len(passports))
  fmt.Printf("%d / %d passports are valid\n", nvalid, len(passports))

  return nil
}
