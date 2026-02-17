---
name: browsing_workflow
description: Learn how to effectively use the browser tool to navigate and interact with web pages.
---

# Browser Automation Workflow

When using the `browser` tool, follow this **Observe-Orient-Decide-Act** loop to avoid errors and "blind" interactions.

## 1. Navigate
Start by navigating to the target URL.
```json
{ "action": "navigate", "url": "https://example.com" }
```

## 2. Observe (CRITICAL STEP)
Do NOT guess selectors (like `#username`, `#login`). Always inspect the page first.
-   **Text-based**: use `get_html` to retrieve the page source.
-   **Visual**: use `screenshot` (if you can process images) to see the layout.

```json
{ "action": "get_html" }
```

## 3. Decide
Analyze the HTML/Image to find the **exact** selectors.
-   Search for `input` tags to find `name`, `id`, or `class` attributes.
-   Search for `button` or `a` tags for login/submit actions.
-   *Example*: If you see `<input name="user_login_123" ...>`, use `input[name='user_login_123']`, not `#username`.

## 4. Act
Perform the action using the verified selector.

```json
{ "action": "type", "selector": "input[name='user_login_123']", "text": "myuser" }
```

## Troubleshooting
-   **"Context deadline exceeded"**: The element was not found within the timeout. Check if the selector is correct or if the page is still loading.
-   **"Node is not visible"**: The element is in the DOM but hidden. It might be behind a menu, requires scrolling, or you are targeting the wrong element (e.g. a hidden input instead of the visible one).
