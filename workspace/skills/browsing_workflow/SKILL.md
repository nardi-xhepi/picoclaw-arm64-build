---
name: browsing_workflow
description: Learn how to effectively use the browser tool to navigate and interact with web pages.
---

# Browser Automation Workflow

## CRITICAL RULE: NAVIGATE FIRST
**If the user provides a URL, your VERY FIRST action in the session MUST be `navigate`.**
Do NOT call `get_html` or `screenshot` until you have successfully navigated to the target URL.

## Correct Flow

1.  **Navigate** (Turn 1)
    -   Action: `{"action": "navigate", "url": "..."}`
    -   **STOP**. Wait for the tool output.

2.  **Observe** (Turn 2)
    -   Action: `{"action": "get_html"}`
    -   **STOP**. Read the HTML to find the *actual* selectors (e.g., `input[name='username']`).
    -   Do NOT guess `#username` or `#login`.

3.  **Act** (Turn 3)
    -   Action: `{"action": "type", "selector": "...", "text": "..."}`

## Troubleshooting
-   **"Page is empty (about:blank)..."**: You forgot to navigate! Call `navigate` immediately.
-   **"You must use 'get_html' or 'screenshot'..."**: You tried to click/type without observing. Call `get_html`.
