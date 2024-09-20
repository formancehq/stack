package hooks

import (
	"bytes"
	"io"
	"net/http"
	"os"
)

type ledgerHook struct {}

func (l ledgerHook) AfterSuccess(hookCtx AfterSuccessContext, res *http.Response) (*http.Response, error) {
	switch hookCtx.HookContext.OperationID {
	case "v2ExportLogs":
		return l.handleV2ExportLogs(hookCtx, res)
	default:
		return res, nil
	}
}

func (l ledgerHook) handleV2ExportLogs(ctx AfterSuccessContext, res *http.Response) (*http.Response, error) {
	if path := ctx.Context.Value("path"); path != nil {
		f, err := os.Create(path.(string))
		if err != nil {
			return nil, err
		}
		_, err = io.Copy(f, res.Body)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

func (l ledgerHook) handleV2ImportLogs(req *http.Request) (*http.Request, error) {
	header := make([]byte, 5)
	_, err := req.Body.Read(header)
	if err == nil && string(header) == "file:" {
		filePath, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}

		f, err := os.OpenFile(string(filePath), os.O_RDONLY, 0666)
		if err != nil {
			return nil, err
		}
		stat, err := f.Stat()
		if err != nil {
			return nil, err
		}

		req.ContentLength = stat.Size()
		req.Body = f

		return req, nil
	}

	data, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	req.Body = io.NopCloser(bytes.NewReader(append(header, data...)))

	return req, nil
}

func (l ledgerHook) BeforeRequest(hookCtx BeforeRequestContext, req *http.Request) (*http.Request, error) {
	switch hookCtx.HookContext.OperationID {
	case "v2ImportLogs":
		return l.handleV2ImportLogs(req)
	default:
		return req, nil
	}
}

