package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

var (
	srcFlag InputSourceFlag
)

func main() {
	flag.Var(&srcFlag, "src",
		"source for the grammar definition; "+
			"empty means use default value; "+
			"\"-\" means use standard input; "+
			"interpreted as filename otherwise",
	)
	flag.Parse()

	r := srcFlag.Get()

	switch rc := r.(type) {
	case io.ReadCloser:
		defer rc.Close()
	}

	defText, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatalf("error reading grammar def: %s", err)
	}

	def, err := parseDefinition(defText)
	if err != nil {
		log.Fatalf("error parsing grammar def: %s", err)
	}

	g, err := def.Build()
	if err != nil {
		log.Fatalf("error building grammar: %s", err)
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	g.Generate(os.Stdout, rng)
}

type InputSourceFlag struct {
	val string
	r   io.Reader
}

func (src *InputSourceFlag) Set(val string) error {
	src.val = val
	if val == "" {
		src.r = strings.NewReader(defaultGrammar)
		return nil
	} else if val == "-" {
		src.r = os.Stdin
		return nil
	} else {
		var err error
		src.r, err = os.Open(val)
		return err
	}
}

func (src *InputSourceFlag) String() string {
	return src.val
}

func (src *InputSourceFlag) Get() io.Reader {
	if src.val == "" && src.r == nil {
		src.r = strings.NewReader(defaultGrammar)
	}
	return src.r
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
