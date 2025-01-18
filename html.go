package prettydiff

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
)

const (
	CSS_STYLES = `* {
        margin: 0;
        padding: 0;
        box-sizing: border-box;
      }

      body {
        background-color: #0d1117;
        color: #c9d1d9;
        font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto,
          "Helvetica Neue", Arial, sans-serif;
      }

      .container {
        width: 100%;
        height: 100vh;
        padding: 25px;
        display: flex;
        flex-direction: column;
        gap: 10px;
      }

      .card {
        border: 1px solid #30363d;
        border-radius: 5px;
      }

      .header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        background-color: #161b22;
        padding: 10px;
        border-bottom: 1px solid #30363d;
      }

      .changes {
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 12px;
        font-size: 14px;
      }

      .add {
        color: #28a745;
      }

      .remove {
        color: #d73a49;
      }

      .body {
        font-size: 14px;
        line-height: 1.5;
      }

      .line {
        display: flex;
        align-items: center;
        padding: 2px 10px;
      }

      .line_number {
        width: 40px;
        text-align: right;
        color: #6a737d;
      }

      .line_content {
        flex: 1;
        padding-left: 10px;
      }

      .added {
        background-color: rgba(50, 207, 94, 0.271);
        color: white;
      }

      .removed {
        background-color: rgba(209, 72, 88, 0.298);
        color: white;
      }

      .file_name {
        font-size: 14px;
      }

      .diff_header {
        background-color: rgba(35, 63, 92, 0.5);
        color: #979ea7;
        padding: 2px 10px;
        font-family: ui-monospace, SFMono-Regular, "SF Mono", Menlo, Consolas, "Liberation Mono", monospace;
        font-size: 12px;
      }

      .large-content {
        display: none;
      }

      .show-content-btn {
        background: rgba(110, 118, 129, 0.4);
        border: none;
        color: #8b949e;
        padding: 4px 8px;
        border-radius: 4px;
        cursor: pointer;
        font-size: 12px;
        margin: 4px 0;
      }

      .show-content-btn:hover {
        background: rgba(110, 118, 129, 0.5);
        color: #c9d1d9;
      }`

	HTML_TEMPLATE = `<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>DIFF</title>
    <style>{{.Styles}}</style>
    <script>
      function toggleContent(id) {
        const content = document.getElementById(id);
        const btn = document.getElementById(id + '-btn');
        if (content.style.display === 'none') {
          content.style.display = 'block';
          btn.style.display = 'none';
        }
      }
    </script>
  </head>
  <body>
    <div class="container">
      {{range .Diffs}}
      <div class="card">
        <div class="header">
          <p class="file_name">{{.FileName}}</p>
          <div class="changes">
            <p class="add">+{{.AddedCount}}</p>
            <p class="remove">-{{.RemovedCount}}</p>
          </div>
        </div>
        {{range .Hunks}}
        <div class="body">
          <div class="diff_header">@@ -{{.OldStart}},{{.OldCount}} +{{.NewStart}},{{.NewCount}} @@{{if .Context}} {{.Context}}{{end}}</div>
          {{range .Lines}}
          <div class="line{{if .Class}} {{.Class}}{{end}}">
            <div class="line_number">{{.LineNumber}}</div>
            <div class="line_content">
              {{if gt (len .Content) 4000}}
                <span id="content-{{.LineNumber}}" class="large-content" style="display: none;">{{.Content}}</span>
                <button id="content-{{.LineNumber}}-btn" class="show-content-btn" onclick="toggleContent('content-{{.LineNumber}}')">
                  Large diff hidden. Click to show {{len .Content}} characters...
                </button>
              {{else}}
                {{.Content}}
              {{end}}
            </div>
          </div>
          {{end}}
        </div>
        {{end}}
      </div>
      {{end}}
    </div>
  </body>
</html>`
)

type HTMLLine struct {
	LineNumber string
	Content    string
	Class      string
}

type HTMLHunk struct {
	OldStart int
	OldCount int
	NewStart int
	NewCount int
	Context  string
	Lines    []HTMLLine
}

type HTMLDiff struct {
	FileName     string
	AddedCount   int
	RemovedCount int
	Hunks        []HTMLHunk
}

type HTMLData struct {
	Styles template.CSS
	Diffs  []HTMLDiff
}

func GenerateHTML(diffs ...*FileDiff) (string, error) {
	var htmlDiffs []HTMLDiff

	for _, diff := range diffs {
		htmlDiff := HTMLDiff{
			FileName: diff.NewFile,
		}

		for _, hunk := range diff.Hunks {
			htmlHunk := HTMLHunk{
				OldStart: hunk.OldStart,
				OldCount: hunk.OldCount,
				NewStart: hunk.NewStart,
				NewCount: hunk.NewCount,
				Context:  hunk.Context,
			}

			for _, change := range hunk.Changes {
				line := HTMLLine{
					Content: change.Content,
				}

				switch change.Type {
				case Unchanged:
					line.LineNumber = fmt.Sprintf("%d", change.NewLine)
				case Added:
					line.LineNumber = fmt.Sprintf("%d", change.NewLine)
					line.Class = "added"
					htmlDiff.AddedCount++
				case Removed:
					line.LineNumber = fmt.Sprintf("%d", change.OldLine)
					line.Class = "removed"
					htmlDiff.RemovedCount++
				}

				htmlHunk.Lines = append(htmlHunk.Lines, line)
			}

			htmlDiff.Hunks = append(htmlDiff.Hunks, htmlHunk)
		}

		htmlDiffs = append(htmlDiffs, htmlDiff)
	}

	data := HTMLData{
		Styles: template.CSS(CSS_STYLES),
		Diffs:  htmlDiffs,
	}

	tmpl, err := template.New("diff").Parse(HTML_TEMPLATE)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %v", err)
	}

	var output strings.Builder
	if err := tmpl.Execute(&output, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %v", err)
	}

	// Minify the output HTML
	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("text/html", html.Minify)

	minified, err := m.String("text/html", output.String())
	if err != nil {
		return "", fmt.Errorf("failed to minify HTML: %v", err)
	}

	return minified, nil
}
