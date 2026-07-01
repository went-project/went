package utils

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

func LoadIgnorePatterns(files ...string) ([]string, string, error) {
	for _, file := range files {
		patterns, err := loadIgnoreList(file)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, "", err
		}
		return patterns, file, nil
	}

	return nil, "", nil
}

func loadIgnoreList(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var patterns []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		patterns = append(patterns, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return patterns, nil
}

func LoadPortFromDotEnv(filenames ...string) (string, string, error) {
	for _, filename := range filenames {
		if _, err := os.Stat(filename); err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return "", "", fmt.Errorf("failed to stat %s: %w", filename, err)
		}

		env, err := godotenv.Read(filename)
		if err != nil {
			return "", "", fmt.Errorf("failed to parse %s: %w", filename, err)
		}
		if port, ok := env["PORT"]; ok {
			return strings.TrimSpace(port), filename, nil
		}
	}

	if port, ok := os.LookupEnv("PORT"); ok {
		return strings.TrimSpace(port), "environment", nil
	}

	return "", "", nil
}

func ShouldIgnore(path string, ignorePatterns []string, root string) bool {
	if len(ignorePatterns) == 0 {
		return false
	}

	rel, err := filepath.Rel(root, path)
	if err != nil {
		rel = path
	}
	rel = normalizePath(rel)

	for _, rawPattern := range ignorePatterns {
		pattern := normalizePath(rawPattern)
		if pattern == "" {
			continue
		}

		if strings.HasSuffix(pattern, "/") {
			trimmed := strings.TrimRight(pattern, "/")
			if rel == trimmed || strings.Contains("/"+rel+"/", "/"+trimmed+"/") {
				return true
			}
			continue
		}

		if rel == pattern || strings.HasPrefix(rel, pattern+"/") || strings.Contains("/"+rel+"/", "/"+pattern+"/") {
			return true
		}

		if match, _ := filepath.Match(pattern, filepath.Base(rel)); match {
			return true
		}
	}

	return false
}

func normalizePath(path string) string {
	cleaned := filepath.ToSlash(filepath.Clean(path))
	if strings.HasPrefix(cleaned, "./") {
		cleaned = strings.TrimPrefix(cleaned, "./")
	}
	cleaned = strings.TrimPrefix(cleaned, "/")
	return cleaned
}
