package events

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/xeipuuv/gojsonschema"
)

func TestGatewayAuditEventSchemas(t *testing.T) {
	for _, name := range []string{"http", "grpc_ok", "grpc_error", "grpc_gateway_error"} {
		t.Run(name, func(t *testing.T) {
			data, err := os.ReadFile(filepath.Join("testdata", "gateway", name+".json"))
			if err != nil {
				t.Fatal(err)
			}
			if err := Check(data, "gateway", "AUDIT"); err != nil {
				t.Fatalf("latest embedded schema: %v", err)
			}
			if err := CheckForVersion(data, "gateway", "v2.0.0", "AUDIT"); err != nil {
				t.Fatalf("versioned embedded schema: %v", err)
			}
			generated, err := os.ReadFile("generated/gateway/v2.0.0/AUDIT.json")
			if err != nil {
				t.Fatal(err)
			}
			schema, err := gojsonschema.NewSchema(gojsonschema.NewBytesLoader(generated))
			if err != nil {
				t.Fatal(err)
			}
			if err := validate(data, schema); err != nil {
				t.Fatalf("generated service schema: %v", err)
			}
			for _, catalog := range []string{"generated/all.json", "../libs/events/generated/all.json"} {
				content, err := os.ReadFile(catalog)
				if err != nil {
					t.Fatal(err)
				}
				var schemas map[string]map[string]map[string]json.RawMessage
				if err := json.Unmarshal(content, &schemas); err != nil {
					t.Fatal(err)
				}
				schema, err := gojsonschema.NewSchema(gojsonschema.NewBytesLoader(schemas["gateway"]["v2.0.0"]["AUDIT"]))
				if err != nil {
					t.Fatal(err)
				}
				if err := validate(data, schema); err != nil {
					t.Fatalf("%s: %v", catalog, err)
				}
			}
		})
	}
}

func TestGatewayAuditRejectsInvalidEvents(t *testing.T) {
	data, err := os.ReadFile("testdata/gateway/grpc_ok.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct {
		name   string
		mutate func(map[string]any)
	}{
		{"wrong envelope type", func(e map[string]any) { e["type"] = "OTHER" }},
		{"wrong app", func(e map[string]any) { e["app"] = "ledger" }},
		{"wrong version", func(e map[string]any) { e["version"] = "v1" }},
		{"missing actor", func(e map[string]any) { delete(e["payload"].(map[string]any), "actor") }},
		{"invalid grpc code", func(e map[string]any) { auditGRPCObject(e)["status_code"] = 17 }},
		{"unknown grpc status is omitted not null", func(e map[string]any) { auditGRPCObject(e)["status_code"] = nil }},
		{"completed requires observed status", func(e map[string]any) { delete(auditGRPCObject(e), "status_code") }},
		{"invalid termination", func(e map[string]any) { auditGRPCObject(e)["termination"] = "success" }},
		{"missing request capture", func(e map[string]any) { delete(auditGRPCObject(e), "request") }},
		{"invalid binary encoding", func(e map[string]any) { auditGRPCObject(e)["request"].(map[string]any)["encoding"] = "hex" }},
		{"invalid base64", func(e map[string]any) { auditGRPCObject(e)["request"].(map[string]any)["data"] = "not base64!" }},
		{"negative byte count", func(e map[string]any) { auditGRPCObject(e)["request"].(map[string]any)["bytes_total"] = -1 }},
		{"capture above 64 KiB", func(e map[string]any) {
			auditGRPCObject(e)["request"].(map[string]any)["data"] = base64.StdEncoding.EncodeToString(bytes.Repeat([]byte{0xff}, 65537))
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			var event map[string]any
			if err := json.Unmarshal(data, &event); err != nil {
				t.Fatal(err)
			}
			test.mutate(event)
			mutated, err := json.Marshal(event)
			if err != nil {
				t.Fatal(err)
			}
			if err := Check(mutated, "gateway", "AUDIT"); err == nil {
				t.Fatal("invalid event accepted")
			}
		})
	}
}

func TestGatewayAuditAcceptsCaptureLimit(t *testing.T) {
	data, err := os.ReadFile("testdata/gateway/grpc_ok.json")
	if err != nil {
		t.Fatal(err)
	}
	var event map[string]any
	if err := json.Unmarshal(data, &event); err != nil {
		t.Fatal(err)
	}
	body := auditGRPCObject(event)["request"].(map[string]any)
	body["data"] = base64.StdEncoding.EncodeToString(bytes.Repeat([]byte{0xff}, 65536))
	body["bytes_total"] = 70000
	body["truncated"] = true
	data, err = json.Marshal(event)
	if err != nil {
		t.Fatal(err)
	}
	if err := Check(data, "gateway", "AUDIT"); err != nil {
		t.Fatalf("valid truncated capture rejected: %v", err)
	}
}

func auditGRPCObject(event map[string]any) map[string]any {
	return event["payload"].(map[string]any)["grpc"].(map[string]any)
}

func TestGatewayAuditAcceptsLegacyHTTPAndAdditiveFields(t *testing.T) {
	data, err := os.ReadFile("testdata/gateway/http.json")
	if err != nil {
		t.Fatal(err)
	}
	var event map[string]any
	if err := json.Unmarshal(data, &event); err != nil {
		t.Fatal(err)
	}
	payload := event["payload"].(map[string]any)
	payload["future_extension"] = map[string]any{"enabled": true}
	actor := payload["actor"].(map[string]any)
	actor["claims"] = map[string]any{"sub": "synthetic-user", "custom_claim": true}
	actor["token_validation_error"] = "synthetic validation error"
	http := payload["http"].(map[string]any)
	http["request"].(map[string]any)["header"] = nil
	http["request"].(map[string]any)["body_truncated"] = true
	http["request"].(map[string]any)["query_params_truncated"] = true
	http["response"].(map[string]any)["headers"] = map[string]any{"X-Empty": nil}
	data, err = json.Marshal(event)
	if err != nil {
		t.Fatal(err)
	}
	if err := Check(data, "gateway", "AUDIT"); err != nil {
		t.Fatalf("compatible HTTP event rejected: %v", err)
	}
}
