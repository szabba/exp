package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
)

type Definition struct {
	entry string
	rules map[string]Rule
}

func parseDefinition(text []byte) (Definition, error) {
	var (
		def     Definition
		rootMap map[string]json.RawMessage
		ruleMap map[string]json.RawMessage
	)

	err := json.Unmarshal(text, &rootMap)
	if err != nil {
		return def, fmt.Errorf("a grammar definition must be a JSON object: %s", err)
	}

	entry, present := rootMap["entry"]
	if !present {
		return def, errors.New("missing grammar entry point")
	}

	rules, present := rootMap["rules"]
	if !present {
		return def, errors.New("missing rules map")
	}

	err = json.Unmarshal(entry, &def.entry)
	if err != nil {
		return def, fmt.Errorf("invalid entry: %s", err)
	}

	err = json.Unmarshal(rules, &ruleMap)
	if err != nil {
		return def, fmt.Errorf("invalid rules map: %s", err)
	}

	def.rules = map[string]Rule{}

	for name, ruleDef := range ruleMap {
		rule, err := parseRule(ruleDef)
		if err != nil {
			log.Printf("%q: %s -> %s", name, ruleDef, err)
			return def, fmt.Errorf("issue parsing rule %q: %s", name, err)
		}
		def.rules[name] = rule
	}

	return def, nil
}

func (def *Definition) Build() (*Grammar, error) {
	var bd Builder
	for name, rule := range def.rules {
		bd.Add(name, rule)
	}
	return bd.Build(def.entry)
}

func parseRule(text []byte) (Rule, error) {
	var errors []error

	lit, err := tryParseLiteral(text)
	if err == nil {
		return lit, err
	}
	errors = append(errors, err)

	ref, err := tryParseReference(text)
	if err == nil {
		return ref, err
	}
	errors = append(errors, err)

	seq, err := tryParseSequence(text)
	if err == nil {
		return seq, err
	}
	errors = append(errors, err)

	alt, err := tryParseAlternative(text)
	if err == nil {
		return alt, err
	}
	errors = append(errors, err)

	return nil, fmt.Errorf("rule %s invalid: %q", text, errors)
}

func tryParseLiteral(text []byte) (Literal, error) {
	var s string
	err := json.Unmarshal(text, &s)
	if err != nil {
		return Literal(s), fmt.Errorf("not a literal -- invalid shape: %s", err)
	}
	return Literal(s), nil
}

func tryParseReference(text []byte) (Reference, error) {
	var (
		parts []string
		err   error
	)

	err = json.Unmarshal(text, &parts)
	if err != nil {
		return Reference(""), fmt.Errorf("invalid reference -- invalid shape: %s", err)
	}

	if len(parts) != 2 {
		return Reference(""), fmt.Errorf("not a reference -- length %d not 2", len(parts))
	}

	if parts[0] != "ref" {
		return Reference(""), fmt.Errorf("not a reference -- wrong prefix %q", parts[0])
	}

	return Reference(parts[1]), err
}

func tryParseSequence(text []byte) (Sequence, error) {
	var (
		seq   Sequence
		err   error
		parts []json.RawMessage
	)

	err = json.Unmarshal(text, &parts)
	if err != nil {
		return nil, fmt.Errorf("not an sequence -- invalid shape: %s", err)
	}

	if len(parts) < 1 {
		return nil, fmt.Errorf("not a sequence -- length %d (< 1)", len(parts))
	}

	if string(parts[0]) != `"seq"` {
		return nil, fmt.Errorf("not a sequence -- wrong prefix %q", parts[0])
	}

	for i, part := range parts[1:] {
		rule, suberr := parseRule(part)
		if suberr != nil {
			return nil, fmt.Errorf("subsequence %d invalid: %s", i, suberr)
		}
		seq = append(seq, rule)
	}

	return seq, err
}

func tryParseAlternative(text []byte) (Alternative, error) {
	var (
		alt   Alternative
		err   error
		parts []json.RawMessage
	)

	err = json.Unmarshal(text, &parts)
	if err != nil {
		return nil, fmt.Errorf("not an alternative -- invalid shape: %s", err)
	}

	if len(parts) < 1 {
		return nil, fmt.Errorf("not an alternative -- length %d (< 2)", len(parts))
	}

	if string(parts[0]) != `"alt"` {
		return nil, fmt.Errorf("not an alternative -- wrong prefix %q", parts[0])
	}

	for i, part := range parts[1:] {
		rule, suberr := parseRule(part)
		if suberr != nil {
			return nil, fmt.Errorf("alternative %d invalid: %s", i, suberr)
		}
		alt = append(alt, rule)
	}
	return alt, err
}
