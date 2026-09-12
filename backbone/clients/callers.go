package backbone

import (
	"context"
	"encoding/json"
	"vague/backbone/process"
)

// Input sends keys in Vim notation.
func (c *Client) Input(keys string) error {
	return c.notify(process.MethodInput, process.InputParams{Keys: keys})
}

// Execute runs a named editor command and decodes its result into result,
// which may be nil to discard it.
func (c *Client) Execute(ctx context.Context, params process.ExecuteParams, result any) error {
	return c.call(ctx, process.MethodExecute, params, result)
}

// Raw runs a command and returns its result as undecoded JSON, for callers
// that just want to display whatever came back.
func (c *Client) Raw(ctx context.Context, params process.ExecuteParams) (json.RawMessage, error) {
	var result json.RawMessage

	if err := c.Execute(ctx, params, &result); err != nil {
		return nil, err
	}

	return result, nil
}
