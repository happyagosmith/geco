package format

import (
	"bytes"
	"fmt"
	"strings"

	yamlv3 "gopkg.in/yaml.v3"
)

type Yaml struct {
	node *yamlv3.Node
}

type Param struct {
	Name         string
	Description  string
	DefaultValue string
	SubSection   string
}

type Section struct {
	Title       string
	Description string
	Params      []Param
}

func NewYaml(b []byte) (Yaml, error) {
	var yn yamlv3.Node
	if err := yamlv3.Unmarshal(b, &yn); err != nil {
		return Yaml{node: nil}, err
	}
	return Yaml{node: &yn}, nil
}

func (y *Yaml) ParseComments() ([]Section, error) {
	c := y.node.Content[0]
	ss := make([]Section, 0, len(c.Content)/2)
	for i := 0; i < len(c.Content)/2; i++ {
		hComment := c.Content[i*2].HeadComment
		if hComment == "" {
			continue
		}

		params := parseParams(c.Content[i*2].Value, "", c.Content[i*2].LineComment, c.Content[i*2+1])
		if len(params) == 0 {
			continue
		}

		cSection := Section{Params: params}

		commentLines := strings.Split(hComment, "#")
		if len(commentLines) > 1 {
			cSection.Title = strings.TrimSuffix(commentLines[1], "\n")
		}

		if len(commentLines) > 2 {
			cSection.Description = strings.Join(commentLines[2:], "")
		}

		ss = append(ss, cSection)
	}

	return ss, nil
}

func parseParams(prefix string, subSection string, lineComment string, n *yamlv3.Node) []Param {
	if len(n.Content) == 0 {
		if len(n.LineComment) == 0 {
			return nil
		}
		dv := n.Value
		if n.Kind == yamlv3.SequenceNode {
			dv = "[]"
		}
		if n.Kind == yamlv3.MappingNode {
			dv = "{}"
		}
		param := Param{
			Name:         prefix,
			Description:  strings.TrimPrefix(n.LineComment, "#"),
			DefaultValue: dv,
			SubSection:   strings.TrimPrefix(subSection, "#")}
		return []Param{param}
	}

	if lineComment != "" {
		dvb, _ := yamlv3.Marshal(n)
		dvs := string(dvb)
		dvs = strings.TrimSuffix(dvs, "\n")
		param := Param{
			Name:         prefix,
			Description:  strings.TrimPrefix(lineComment, "#"),
			DefaultValue: dvs,
			SubSection:   strings.TrimPrefix(subSection, "#")}

		return []Param{param}
	}

	if n.Kind == yamlv3.MappingNode {
		ps := make([]Param, 0, len(n.Content)/2)
		for i := 0; i < len(n.Content)/2; i++ {
			nPrefix := strings.Join([]string{prefix, n.Content[i*2].Value}, ".")
			ss := subSection
			if hc := n.Content[i*2].HeadComment; hc != "" {
				ss = hc
			}
			ps = append(ps, parseParams(nPrefix, ss, n.Content[i*2].LineComment, n.Content[i*2+1])...)
		}
		return ps
	}

	return []Param{}
}

func (y *Yaml) Merge(oy Yaml) error {
	err := mergeNodes(y.node, oy.node)
	return err
}

func (y Yaml) String() (string, error) {
	out, err := y.Bytes()
	if err != nil {
		return "", err
	}

	return strings.TrimSuffix(string(out), "\n"), nil
}

func (y Yaml) Bytes() ([]byte, error) {
	var out bytes.Buffer
	encoder := yamlv3.NewEncoder(&out)
	encoder.SetIndent(1)
	err := encoder.Encode(y.node)
	if err != nil {
		return nil, err
	}
	defer encoder.Close()

	return out.Bytes(), nil
}

func mergeNodes(a, b *yamlv3.Node) error {
	if a.Kind != b.Kind {
		return fmt.Errorf("it is not possible to merge different types")
	}

	if a.Kind == yamlv3.DocumentNode {
		err := mergeNodes(a.Content[0], b.Content[0])
		if err != nil {
			return err
		}
	}

	if a.Kind == yamlv3.MappingNode {
		lmb := lookUpMap(b.Content)

		for i := 0; i < len(a.Content)/2; i++ {
			key := a.Content[i*2].Value
			if n, ok := lmb[key]; ok {
				if n.nodeValue.Kind == yamlv3.MappingNode {
					_ = mergeNodes(a.Content[i*2+1], n.nodeValue)
					a.Column = a.Column - 2
				} else {
					a.Content[i*2+1] = n.nodeValue
				}
				n.found = true
				lmb[key] = n
			}
		}

		appendContent(a, b, lmb)
	}

	return nil
}

func ToYaml(v any, indent int) string {
	data, err := yamlv3.Marshal(v)
	if err != nil {
		return ""
	}
	is := "\n" + strings.Repeat(" ", indent)
	ys := string(data)
	ys = is + strings.ReplaceAll(ys, "\n", is)
	ys = strings.TrimSuffix(ys, is)
	return ys
}

type node struct {
	nodeKey   *yamlv3.Node
	nodeValue *yamlv3.Node
	found     bool
}

func lookUpMap(nodes []*yamlv3.Node) map[string]node {
	nb := map[string]node{}
	for i := 0; i < len(nodes)/2; i++ {
		key := nodes[i*2]
		value := nodes[i*2+1]
		nb[key.Value] = node{
			nodeKey:   key,
			nodeValue: value,
			found:     false,
		}
	}

	return nb
}

func appendContent(a *yamlv3.Node, b *yamlv3.Node, lmb map[string]node) {
	for i := 0; i < len(b.Content)/2; i++ {
		key := b.Content[i*2].Value
		if n := lmb[key]; !n.found {
			a.Content = append(a.Content, lmb[key].nodeKey)
			a.Content = append(a.Content, lmb[key].nodeValue)
		}
	}
}
