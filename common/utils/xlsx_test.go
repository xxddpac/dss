package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestXlsx(t *testing.T) {
	var (
		name   = fmt.Sprintf("test_scan_result_%s.xlsx", time.Now().Format(TimeLayout))
		file   = filepath.Join(os.TempDir(), name)
		result [][]string
	)
	defer func() {
		if err := os.Remove(file); err != nil {
			t.Fatal(err)
		}
	}()
	result = append(result,
		[]string{"8.8.8.8", "22", time.Now().Format(TimeLayout)},
		[]string{"9.9.9.9", "23", time.Now().Format(TimeLayout)},
	)
	if err := WriteToXlsx(file, []string{"host", "port", "date"}, result); err != nil {
		t.Fatal(err)
	}
}
