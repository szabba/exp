package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"time"
)

func main() {
	def, err := parseDefinition([]byte(defaultGrammar))
	if err != nil {
		log.Fatalf("error parsing default grammar: %s", err)
	}

	g, _ := def.Build()
	rng := rand.New(rand.NewSource(time.Now().Unix()))
	g.Generate(os.Stdout, rng)
}

type Builder struct {
	rules   map[string]Rule
	missing map[string]int
}

func (bd *Builder) init() {
	if bd.rules != nil && bd.missing != nil {
		return
	}
	bd.rules = map[string]Rule{}
	bd.missing = map[string]int{}
}

func (bd *Builder) Add(name string, rule Rule) error {
	bd.init()

	if _, present := bd.rules[name]; present {
		return fmt.Errorf("trying to add a duplicate rule for name %q", name)
	}

	bd.rules[name] = rule
	rule.EachRef(bd.ruleShouldExist)
	delete(bd.missing, name)
	return nil
}

func (bd *Builder) ruleShouldExist(name string) {
	if _, present := bd.rules[name]; !present {
		bd.missing[name]++
	}
}

func (bd *Builder) MissingRules() []string {
	names := make([]string, 0, len(bd.missing))
	for name := range bd.missing {
		names = append(names, name)
	}
	return names
}

func (bd *Builder) Build(entry string) (*Grammar, error) {
	if bd.hasMissingRules() {
		return nil, fmt.Errorf("missing rules for names %q", bd.MissingRules())
	}
	if bd.noRuleForName(entry) {
		return nil, fmt.Errorf("missing rule for entry point %q", entry)
	}
	defer bd.init()
	return &Grammar{entry, bd.rules}, nil
}

func (bd *Builder) hasMissingRules() bool {
	return len(bd.missing) > 0
}

func (bd *Builder) noRuleForName(name string) bool {
	_, present := bd.rules[name]
	_, missing := bd.missing[name]
	return !present && missing
}

type Grammar struct {
	entry string
	rules map[string]Rule
}

func (g *Grammar) Generate(into io.Writer, rng *rand.Rand) error {
	return g.rules[g.entry].GenerateOne(into, g.rules, rng)
}

type Rule interface {
	GenerateOne(into io.Writer, others map[string]Rule, rng *rand.Rand) error
	EachRef(f func(string))
}

type Alternative []Rule

var _ Rule = Alternative{}

func (alt Alternative) GenerateOne(w io.Writer, rules map[string]Rule, rng *rand.Rand) error {
	if len(alt) == 0 {
		return errors.New("empty alternative")
	}
	rule := alt[rng.Intn(len(alt))]
	return rule.GenerateOne(w, rules, rng)
}

func (alt Alternative) EachRef(f func(string)) {
	for _, rule := range alt {
		rule.EachRef(f)
	}
}

type Sequence []Rule

var _ Rule = Sequence{}

func (seq Sequence) GenerateOne(w io.Writer, rules map[string]Rule, rng *rand.Rand) error {
	for _, rule := range seq {
		if err := rule.GenerateOne(w, rules, rng); err != nil {
			return err
		}
	}
	return nil
}

func (seq Sequence) EachRef(f func(string)) {
	for _, rule := range seq {
		rule.EachRef(f)
	}
}

type Reference string

var _ Rule = Reference("")

func (ref Reference) GenerateOne(w io.Writer, rules map[string]Rule, rng *rand.Rand) error {
	rule, ok := rules[string(ref)]
	if !ok {
		return fmt.Errorf("no rule named %q", ref)
	}
	return rule.GenerateOne(w, rules, rng)
}

func (ref Reference) EachRef(f func(string)) {
	f(string(ref))
}

type Literal string

var _ Rule = Literal("")

func (l Literal) GenerateOne(w io.Writer, _ map[string]Rule, _ *rand.Rand) error {
	_, err := io.WriteString(w, string(l))
	return err
}

func (l Literal) EachRef(_ func(string)) {}
