package enrichment

import (
	"errors"
	"fmt"
	"github.com/jitsucom/eventnative/jsonutils"
	"strings"
)

type Rule interface {
	Execute(fact map[string]interface{}) error
}

func NewRule(ruleConfig *RuleConfig) (Rule, error) {
	err := ruleConfig.Validate()
	if err != nil {
		return nil, err
	}

	source := jsonutils.NewJsonPath(ruleConfig.From)
	if source.IsEmpty() {
		return nil, errors.New("'from' must be a valid path like: /node1/node2")
	}
	destination := jsonutils.NewJsonPath(ruleConfig.To)
	if destination.IsEmpty() {
		return nil, errors.New("'to' must be a valid path like: /node1/node2")
	}

	switch ruleConfig.Name {
	case IpLookup:
		return NewIpLookupRule(source, destination)
	case UserAgentParse:
		return NewUserAgentParseRule(source, destination)
	default:
		return nil, fmt.Errorf("Unsupported enrichment rule type: %s", ruleConfig.Name)
	}
}

type RuleConfig struct {
	Name string `mapstructure:"name" json:"name,omitempty" yaml:"name,omitempty"`
	From string `mapstructure:"from" json:"from,omitempty" yaml:"from,omitempty"`
	To   string `mapstructure:"to" json:"to,omitempty" yaml:"to,omitempty"`
}

func (r *RuleConfig) Validate() error {
	r.Name = strings.ToLower(r.Name)
	r.To = strings.ToLower(r.To)
	r.From = strings.ToLower(r.From)

	if r.To == "" {
		return errors.New("'to' is required enrichment rule parameter")
	}

	if r.From == "" {
		return errors.New("'from' is required enrichment rule parameter")
	}

	return nil
}
