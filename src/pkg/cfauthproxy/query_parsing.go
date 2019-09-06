package cfauthproxy

import (
	"errors"

	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/promql"
)

type QueryParser struct{}

func (q *QueryParser) ExtractSourceIds(query string) ([]string, error) {
	expr, err := promql.ParseExpr(query)
	if err != nil {
		return nil, err
	}

	visitor := newSourceIdVisitor()

	err = promql.Walk(
		visitor,
		expr,
		nil,
	)
	if err != nil {
		return nil, err
	}

	var sourceIds []string

	for sourceId := range visitor.sourceIds {
		sourceIds = append(sourceIds, sourceId)
	}

	return sourceIds, nil
}

type sourceIdVisitor struct {
	sourceIds map[string]struct{}
}

func newSourceIdVisitor() *sourceIdVisitor {
	return &sourceIdVisitor{
		sourceIds: make(map[string]struct{}),
	}
}

func (s *sourceIdVisitor) Visit(node promql.Node, _ []promql.Node) (promql.Visitor, error) {
	if node == nil {
		return nil, nil
	}

	var err error
	switch selector := node.(type) {
	case *promql.VectorSelector:
		err = s.addSourceIdsFromMatchers(selector.LabelMatchers)
	case *promql.MatrixSelector:
		err = s.addSourceIdsFromMatchers(selector.LabelMatchers)
	}

	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *sourceIdVisitor) addSourceIdsFromMatchers(labelMatchers []*labels.Matcher) error {
	containsValidSourceId := false
	for _, labelMatcher := range labelMatchers {
		if labelMatcher.Name == "source_id" && labelMatcher.Value != "" {
			containsValidSourceId = true
			err := addSourceIdsFromLabelMatcher(s.sourceIds, labelMatcher)
			if err != nil {
				return err
			}
		}
	}

	if !containsValidSourceId {
		return errors.New("one or more terms lack a sourceId")
	}

	return nil
}

func addSourceIdsFromLabelMatcher(sourceIds map[string]struct{}, labelMatcher *labels.Matcher) error {
	switch labelMatcher.Type {
	case labels.MatchRegexp:
		return errors.New("regular expressions are unavailable on source ids")
	case labels.MatchEqual:
		sourceIds[labelMatcher.Value] = struct{}{}
	}
	return nil
}
