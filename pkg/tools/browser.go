package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/chromedp/chromedp"
)

type BrowserTool struct {
	userDataDir string
	headless    bool
}

func NewBrowserTool(userDataDir string, headless bool) *BrowserTool {
	return &BrowserTool{
		userDataDir: userDataDir,
		headless:    headless,
	}
}

func (t *BrowserTool) Name() string {
	return "browser"
}

func (t *BrowserTool) Description() string {
	return "Automate a web browser. CRITICAL: Do NOT chain 'navigate' and 'click'/'type' in one turn. You MUST 1) 'navigate', 2) STOP and wait for result, 3) 'get_html' to find selectors, 4) 'click'/'type'."
}

func (t *BrowserTool) Parameters() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"action": map[string]interface{}{
				"type":        "string",
				"description": "Action to perform: navigate, click, type, screenshot, get_html, evaluate",
				"enum":        []string{"navigate", "click", "type", "screenshot", "get_html", "evaluate"},
			},
			"url": map[string]interface{}{
				"type":        "string",
				"description": "URL to navigate to (for action=navigate)",
			},
			"selector": map[string]interface{}{
				"type":        "string",
				"description": "CSS selector for element to interact with (for click, type)",
			},
			"text": map[string]interface{}{
				"type":        "string",
				"description": "Text to type (for action=type)",
			},
			"skip_wait": map[string]interface{}{
				"type":        "boolean",
				"description": "If true, skip waiting for element visibility before interaction (default: false)",
			},
			"script": map[string]interface{}{
				"type":        "string",
				"description": "JavaScript to evaluate (for action=evaluate)",
			},
		},
		"required": []string{"action"},
	}
}

func (t *BrowserTool) Execute(ctx context.Context, args map[string]interface{}) *ToolResult {
	action, ok := args["action"].(string)
	if !ok {
		return ErrorResult("action is required")
	}

	// Prepare Chrome options
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.NoSandbox,
		//chromedp.Flag("disable-gpu", true), // already in DefaultExecAllocatorOptions?
	)

	if t.headless {
		opts = append(opts, chromedp.Headless)
	} else {
		opts = append(opts, chromedp.Flag("headless", false))
	}

	if t.userDataDir != "" {
		absPath, err := filepath.Abs(t.userDataDir)
		if err == nil {
			opts = append(opts, chromedp.UserDataDir(absPath))
		}
	}

	// Create allocator context
	allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	// Create context
	ctx, cancel = chromedp.NewContext(allocCtx)
	defer cancel()

	// Set timeout
	ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	var result string
	var err error

	switch action {
	case "navigate":
		urlStr, ok := args["url"].(string)
		if !ok {
			return ErrorResult("url is required for navigate action")
		}
		err = chromedp.Run(ctx,
			chromedp.Navigate(urlStr),
		)
		if err == nil {
			result = fmt.Sprintf("Navigated to %s", urlStr)
		}

	case "click":
		selector, ok := args["selector"].(string)
		if !ok {
			return ErrorResult("selector is required for click action")
		}
		// Default to waiting for visible, unless skip_wait is true
		waitVisible := true
		if skip, ok := args["skip_wait"].(bool); ok && skip {
			waitVisible = false
		}

		actions := []chromedp.Action{}
		if waitVisible {
			actions = append(actions, chromedp.WaitVisible(selector))
		}
		actions = append(actions, chromedp.Click(selector, chromedp.NodeVisible))

		err = chromedp.Run(ctx, actions...)
		if err == nil {
			result = fmt.Sprintf("Clicked element %s", selector)
		}

	case "type":
		selector, ok := args["selector"].(string)
		if !ok {
			return ErrorResult("selector is required for type action")
		}
		text, ok := args["text"].(string)
		if !ok {
			return ErrorResult("text is required for type action")
		}

		actions := []chromedp.Action{
			chromedp.WaitVisible(selector),
			chromedp.SendKeys(selector, text, chromedp.NodeVisible),
		}

		err = chromedp.Run(ctx, actions...)
		if err == nil {
			result = fmt.Sprintf("Typed '%s' into %s", text, selector)
		}

	case "screenshot":
		var buf []byte
		err = chromedp.Run(ctx,
			chromedp.CaptureScreenshot(&buf),
		)
		if err == nil {
			// Save screenshot to disk? Or return base64?
			// For now, let's return base64 so LLM can potentially see it (if supported) or just ack
			// Ideally we save to file and return path
			// But for tool result, let's just say "Screenshot captured" and maybe save to temp?

			// Actually, let's save to a temp file for user to see
			tmpFile, cleanErr := os.CreateTemp("", "screenshot-*.png")
			if cleanErr == nil {
				defer tmpFile.Close()
				tmpFile.Write(buf)
				result = fmt.Sprintf("Screenshot saved to %s", tmpFile.Name())
			} else {
				result = fmt.Sprintf("Screenshot captured (%d bytes) but failed to save file: %v", len(buf), cleanErr)
			}
		}

	case "get_html":
		var content string
		err = chromedp.Run(ctx,
			chromedp.OuterHTML("html", &content),
		)
		if err == nil {
			result = content
			// Truncate for display
			if len(result) > 5000 {
				result = result[:5000] + "... (truncated)"
			}
		}

	case "evaluate":
		script, ok := args["script"].(string)
		if !ok {
			return ErrorResult("script is required for evaluate action")
		}
		var res interface{}
		err = chromedp.Run(ctx,
			chromedp.Evaluate(script, &res),
		)
		if err == nil {
			jsonRes, _ := json.Marshal(res)
			result = string(jsonRes)
		}

	default:
		return ErrorResult(fmt.Sprintf("unknown action: %s", action))
	}

	if err != nil {
		// Attempt to get current URL and Title for context
		var currentURL, title string
		chromedp.Run(ctx,
			chromedp.Location(&currentURL),
			chromedp.Title(&title),
		)
		return ErrorResult(fmt.Sprintf("Browser action failed: %v\nContext: %s (%s)", err, title, currentURL))
	}

	return &ToolResult{
		ForLLM:  result,
		ForUser: result,
	}
}
