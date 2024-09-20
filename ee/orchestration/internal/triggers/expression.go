package triggers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.temporal.io/sdk/temporal"

	"github.com/formancehq/go-libs/collectionutils"

	"github.com/expr-lang/expr"
	"github.com/formancehq/go-libs/api"
	"github.com/pkg/errors"
)

type expressionEvaluator struct {
	httpClient *http.Client
}

func (h *expressionEvaluator) link(params ...any) (any, error) {
	if len(params) != 2 {
		return nil, temporal.NewNonRetryableApplicationError(
			fmt.Sprintf("expect two arguments, got %d", len(params)),
			"APPLICATION",
			fmt.Errorf("expect two arguments, got %d", len(params)),
		)
	}

	data, _ := json.Marshal(params[0])

	type object struct {
		Links []api.Link `json:"links"`
	}
	o := &object{}
	if err := json.Unmarshal(data, o); err != nil {
		return nil, err
	}

	rel, ok := params[1].(string)
	if !ok {
		return nil, errors.New("second parameter must be a string")
	}

	filteredLinks := collectionutils.Filter(o.Links, func(link api.Link) bool {
		return link.Name == rel
	})

	switch len(filteredLinks) {
	case 0:
		return nil, temporal.NewNonRetryableApplicationError(
			fmt.Sprintf("link '%s' not defined for object", rel),
			"APPLICATION",
			fmt.Errorf("link '%s' not defined for object", rel),
		)
	case 1:
		rsp, err := h.httpClient.Get(filteredLinks[0].URI)
		if err != nil {
			return nil, errors.Wrapf(err, "reading resource: %s", filteredLinks[0].URI)
		}
		if rsp.StatusCode >= 400 {
			return nil, fmt.Errorf("unexpected status code when reading resource: %d", rsp.StatusCode)
		}

		apiResponse := api.BaseResponse[map[string]any]{}
		if err := json.NewDecoder(rsp.Body).Decode(&apiResponse); err != nil {
			return nil, errors.Wrap(err, "decoding response")
		}

		return apiResponse.Data, nil
	default:
		return nil, temporal.NewNonRetryableApplicationError(
			fmt.Sprintf("multiple link '%s' found for object", rel),
			"APPLICATION",
			fmt.Errorf("multiple link '%s' found for object", rel),
		)
	}
}

func (h *expressionEvaluator) eval(rawObject any, e string) (any, error) {
	p, err := expr.Compile(e, expr.Function("link", h.link))
	if err != nil {
		return "", err
	}

	ret, err := expr.Run(p, map[string]any{
		"event": rawObject,
	})
	if err != nil {
		if err := errors.Unwrap(err); err != nil {
			return nil, err
		}
		return nil, err
	}

	return ret, nil
}

func (h *expressionEvaluator) evalFilter(event any, filter string) (bool, error) {
	ret, err := h.eval(event, filter)
	if err != nil {
		return false, err
	}

	switch ret := ret.(type) {
	case bool:
		return ret, nil
	default:
		return false, nil
	}
}

func (h *expressionEvaluator) evalVariable(rawObject any, e string) (string, error) {
	ret, err := h.eval(rawObject, e)
	if err != nil {
		return "", err
	}

	switch ret.(type) {
	case float64, float32:
		data, err := json.Marshal(ret)
		if err != nil {
			return "", err
		}
		return string(data), nil
	default:
		return fmt.Sprint(ret), nil
	}
}

func (h *expressionEvaluator) evalVariables(rawObject any, vars map[string]string) (map[string]string, error) {
	results := make(map[string]string)
	for k, v := range vars {
		var err error
		results[k], err = h.evalVariable(rawObject, v)
		if err != nil {
			return nil, err
		}
	}

	return results, nil
}

func NewExpressionEvaluator(httpClient *http.Client) *expressionEvaluator {
	return &expressionEvaluator{
		httpClient: httpClient,
	}
}

func NewDefaultExpressionEvaluator() *expressionEvaluator {
	return NewExpressionEvaluator(http.DefaultClient)
}
