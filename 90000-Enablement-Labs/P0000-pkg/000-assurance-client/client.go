package assurance

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
)

type BridgeClient struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout io.ReadCloser
	stderr io.ReadCloser
}

func NewBridgeClient(binaryPath string) (*BridgeClient, error) {
	cmd := exec.Command(binaryPath)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	go io.Copy(os.Stderr, stderr)

	return &BridgeClient{
		cmd:    cmd,
		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,
	}, nil
}

func (c *BridgeClient) CallTool(ctx context.Context, toolName string, args map[string]interface{}) (string, error) {
	// Simplified JSON-RPC for MCP stdio
	req := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "tools/call",
		"params": map[string]interface{}{
			"name":      toolName,
			"arguments": args,
		},
	}

	data, _ := json.Marshal(req)
	fmt.Fprintln(c.stdin, string(data))

	// In a real implementation, we'd have a proper JSON-RPC reader loop.
	// For this assurance suite, we'll read the next line.
	var resp struct {
		Result struct {
			Content []struct {
				Text string `json:"Olympus2/00000-Identity-Foundations/P0000-pkg/text"`
			} `json:"content"`
			IsError bool `json:"isError"`
		} `json:"result"`
		Error *struct {
			Message string `json:"message"`
		} `json:"error"`
	}

	dec := json.NewDecoder(c.stdout)
	if err := dec.Decode(&resp); err != nil {
		return "", err
	}

	if resp.Error != nil {
		return "", fmt.Errorf(resp.Error.Message)
	}

	if len(resp.Result.Content) > 0 {
		return resp.Result.Content[0].Text, nil
	}

	return "", fmt.Errorf("no content in response")
}

func (c *BridgeClient) Close() {
	c.stdin.Close()
	c.cmd.Process.Kill()
}
