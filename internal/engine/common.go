package engine

import (
	"fmt"
	"strings"

	"github.com/Azure/InnovationEngine/internal/lib"
	"github.com/Azure/InnovationEngine/internal/logging"
	"github.com/Azure/InnovationEngine/internal/ui"
	"github.com/xrash/smetrics"
)

// Indents a multi-line command to be nested under the first line of the
// command.
func indentMultiLineCommand(content string, indentation int) string {
	lines := strings.Split(content, "\n")
	for i := 1; i < len(lines); i++ {
		if strings.HasSuffix(strings.TrimSpace(lines[i-1]), "\\") {
			lines[i] = strings.Repeat(" ", indentation) + lines[i]
		} else if strings.TrimSpace(lines[i]) != "" {
			lines[i] = strings.Repeat(" ", indentation) + lines[i]
		}

	}
	return strings.Join(lines, "\n")
}

// Compares the actual output of a command to the expected output of a command.
func compareCommandOutputs(
	actualOutput string,
	expectedOutput string,
	expectedSimilarity float64,
	expectedOutputLanguage string,
) error {
	if strings.ToLower(expectedOutputLanguage) == "json" {
		logging.GlobalLogger.Debugf(
			"Comparing JSON strings:\nExpected: %s\nActual%s",
			expectedOutput,
			actualOutput,
		)
		results, err := lib.CompareJsonStrings(actualOutput, expectedOutput, expectedSimilarity)

		if err != nil {
			return err
		}

		if !results.AboveThreshold {
			return fmt.Errorf(
				ui.ErrorMessageStyle.Render("Expected output does not match actual output."),
			)
		}

		logging.GlobalLogger.Debugf(
			"Expected Similarity: %f, Actual Similarity: %f",
			expectedSimilarity,
			results.Score,
		)
	} else {
		score := smetrics.JaroWinkler(expectedOutput, actualOutput, 0.7, 4)

		if expectedSimilarity > score {
			return fmt.Errorf(ui.ErrorMessageStyle.Render("Expected output does not match actual output."))
		}
	}

	return nil
}
