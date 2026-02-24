package backend

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRender(t *testing.T) {
	b := &Backend{}

	tests := []struct {
		name     string
		input    string
		contains []string
	}{
		{
			name: "table with alignment",
			input: `| Left | Center | Right |
|:-----|:------:|------:|
| a    | b      | c     |
`,
			contains: []string{
				"<table>",
				"<thead>",
				"<tbody>",
				`<th style="text-align:left">Left</th>`,
				`<th style="text-align:center">Center</th>`,
				`<th style="text-align:right">Right</th>`,
			},
		},
		{
			name: "footnotes",
			input: `This has a footnote[^1].

[^1]: Footnote content.
`,
			contains: []string{
				"<sup",
				"Footnote content",
			},
		},
		{
			name: "unchecked task list",
			input: `- [ ] unchecked item
`,
			contains: []string{
				`<input`,
				`type="checkbox"`,
				"unchecked item",
			},
		},
		{
			name: "checked task list",
			input: `- [x] checked item
`,
			contains: []string{
				`<input`,
				`checked`,
				"checked item",
			},
		},
		{
			name: "mixed task list",
			input: `- [x] done
- [ ] todo
`,
			contains: []string{
				`checked`,
				"done",
				"todo",
			},
		},
		{
			name:  "strikethrough",
			input: `~~deleted~~`,
			contains: []string{
				"<del>deleted</del>",
			},
		},
		{
			name:  "fenced code block with language",
			input: "```go\nfmt.Println(\"hello\")\n```\n",
			contains: []string{
				`<code class="language-go"`,
				`fmt.Println`,
			},
		},
		{
			name:  "autolinks",
			input: `Visit https://example.com for more.`,
			contains: []string{
				`<a href="https://example.com"`,
				"https://example.com",
			},
		},
		{
			name: "normal list should not become checklist",
			input: `- item one
- item two
`,
			contains: []string{
				"<li>item one</li>",
				"<li>item two</li>",
			},
		},
		{
			name:  "inline math",
			input: `The formula $x^2 + y^2 = z^2$ is well known.`,
			contains: []string{
				`class="katex"`,
				`x`,
				`2`,
			},
		},
		{
			name:  "block math",
			input: "$$E = mc^2$$\n",
			contains: []string{
				`class="katex-display"`,
				`E`,
			},
		},
		{
			name:  "math with subscript preserved",
			input: `The weight $w_k$ of term k.`,
			contains: []string{
				`class="katex"`,
				`w`,
			},
		},
		{
			name: "combined markdown features",
			input: `# Heading

| Feature | Status |
|---------|--------|
| Tables  | ~~old~~ |

- [x] task done
- [ ] task todo

Check https://example.com and this footnote[^1].

[^1]: A note.
`,
			contains: []string{
				`id="heading"`,
				"<table>",
				"<del>old</del>",
				`checked`,
				`type="checkbox"`,
				`<a href="https://example.com"`,
				"A note",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := b.Render([]byte(tt.input))
			assert.NoError(t, err)
			html := string(result)
			for _, s := range tt.contains {
				assert.Contains(t, html, s)
			}
		})
	}
}
