package cmd

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"

	"github.com/tomba-io/tomba/pkg/util"
	"github.com/tomba-io/tomba/pkg/version"
)

const updateRepo = "tomba-io/tomba"

var updateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"upgrade"},
	Short:   "Update tomba to the latest version.",
	Long:    Long,
	Example: updateExample,
	Run:     updateRun,
}

func updateRun(cmd *cobra.Command, args []string) {
	current := version.Version
	fmt.Printf("%s Current version: %s\n", util.InfoIcon(), util.Bold("v"+current))

	// Fetch latest version
	fmt.Printf("%s Checking for updates...\n", util.InfoIcon())
	latest, err := getLatestVersion()
	if err != nil {
		fmt.Printf("%s Failed to check for updates: %s\n", util.ErrorIcon(), util.Red(err.Error()))
		return
	}

	// Compare
	if strings.TrimPrefix(current, "v") == strings.TrimPrefix(latest, "v") {
		fmt.Printf("%s Already up to date (v%s)\n", util.SuccessIcon(), current)
		return
	}

	fmt.Printf("%s New version available: %s -> %s\n", util.InfoIcon(),
		util.Yellow("v"+current), util.Green("v"+latest))

	// Determine archive filename
	osName := runtime.GOOS
	archName := runtime.GOARCH
	ext := "tar.gz"
	if osName == "windows" {
		ext = "zip"
	}
	filename := fmt.Sprintf("tomba_%s_%s.%s", osName, archName, ext)
	downloadURL := fmt.Sprintf("https://github.com/%s/releases/download/v%s/%s", updateRepo, latest, filename)
	checksumsURL := fmt.Sprintf("https://github.com/%s/releases/download/v%s/tomba_checksums.txt", updateRepo, latest)

	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "tomba-update-*")
	if err != nil {
		fmt.Printf("%s Failed to create temp directory: %s\n", util.ErrorIcon(), util.Red(err.Error()))
		return
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Download archive
	fmt.Printf("%s Downloading %s...\n", util.InfoIcon(), util.Cyan(filename))
	archivePath := filepath.Join(tmpDir, filename)
	if err := downloadFile(downloadURL, archivePath); err != nil {
		fmt.Printf("%s Download failed: %s\n", util.ErrorIcon(), util.Red(err.Error()))
		return
	}
	fmt.Printf("%s Downloaded successfully\n", util.SuccessIcon())

	// Download and verify checksum
	checksumsPath := filepath.Join(tmpDir, "checksums.txt")
	if err := downloadFile(checksumsURL, checksumsPath); err == nil {
		if err := verifyChecksum(archivePath, filename, checksumsPath); err != nil {
			fmt.Printf("%s %s\n", util.ErrorIcon(), util.Red(err.Error()))
			return
		}
		fmt.Printf("%s Checksum verified (SHA256)\n", util.SuccessIcon())
	} else {
		fmt.Printf("%s Could not download checksums, skipping verification\n", util.WarningIcon())
	}

	// Extract binary
	fmt.Printf("%s Extracting...\n", util.InfoIcon())
	binaryName := "tomba"
	if osName == "windows" {
		binaryName = "tomba.exe"
	}
	extractedPath := filepath.Join(tmpDir, binaryName)
	if err := extractBinary(archivePath, tmpDir, binaryName, osName); err != nil {
		fmt.Printf("%s Extraction failed: %s\n", util.ErrorIcon(), util.Red(err.Error()))
		return
	}
	fmt.Printf("%s Extracted successfully\n", util.SuccessIcon())

	// Find current binary path
	execPath, err := os.Executable()
	if err != nil {
		fmt.Printf("%s Cannot determine binary path: %s\n", util.ErrorIcon(), util.Red(err.Error()))
		return
	}
	execPath, err = filepath.EvalSymlinks(execPath)
	if err != nil {
		fmt.Printf("%s Cannot resolve binary path: %s\n", util.ErrorIcon(), util.Red(err.Error()))
		return
	}

	// Replace binary
	fmt.Printf("%s Installing to %s...\n", util.InfoIcon(), util.Bold(execPath))
	oldPath := execPath + ".old"

	// Rename current binary
	if err := os.Rename(execPath, oldPath); err != nil {
		fmt.Printf("%s Cannot rename current binary (may need sudo): %s\n", util.ErrorIcon(), util.Red(err.Error()))
		return
	}

	// Move new binary into place
	if err := copyFile(extractedPath, execPath); err != nil {
		// Restore old binary on failure
		_ = os.Rename(oldPath, execPath)
		fmt.Printf("%s Cannot install new binary: %s\n", util.ErrorIcon(), util.Red(err.Error()))
		return
	}

	// Make executable
	if err := os.Chmod(execPath, 0755); err != nil {
		fmt.Printf("%s Cannot set permissions: %s\n", util.WarningIcon(), util.Yellow(err.Error()))
	}

	// Remove old binary
	_ = os.Remove(oldPath)

	fmt.Printf("\n%s Updated tomba %s -> %s\n", util.SuccessIcon(),
		util.Yellow("v"+current), util.Green("v"+latest))
}

func getLatestVersion() (string, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", updateRepo))
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API returned %d", resp.StatusCode)
	}

	var release struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	v := strings.TrimPrefix(release.TagName, "v")
	if v == "" {
		return "", fmt.Errorf("empty version in response")
	}
	return v, nil
}

func downloadFile(url, dest string) error {
	resp, err := http.Get(url) //nolint:gosec
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer func() { _ = out.Close() }()

	_, err = io.Copy(out, resp.Body)
	return err
}

func verifyChecksum(file, filename, checksumsFile string) error {
	data, err := os.ReadFile(checksumsFile)
	if err != nil {
		return fmt.Errorf("cannot read checksums: %w", err)
	}

	var expected string
	for _, line := range strings.Split(string(data), "\n") {
		parts := strings.Fields(line)
		if len(parts) == 2 && parts[1] == filename {
			expected = parts[0]
			break
		}
	}
	if expected == "" {
		return fmt.Errorf("checksum not found for %s", filename)
	}

	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return err
	}
	actual := hex.EncodeToString(h.Sum(nil))

	if actual != expected {
		return fmt.Errorf("checksum mismatch: expected %s, got %s", expected, actual)
	}
	return nil
}

func extractBinary(archivePath, destDir, binaryName, osName string) error {
	if osName == "windows" {
		return extractZip(archivePath, destDir, binaryName)
	}
	return extractTarGz(archivePath, destDir, binaryName)
}

func extractTarGz(archivePath, destDir, binaryName string) error {
	f, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer func() { _ = gz.Close() }()

	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if filepath.Base(hdr.Name) == binaryName && hdr.Typeflag == tar.TypeReg {
			out, err := os.Create(filepath.Join(destDir, binaryName))
			if err != nil {
				return err
			}
			if _, err := io.Copy(out, tr); err != nil { //nolint:gosec
				_ = out.Close()
				return err
			}
			_ = out.Close()
			return nil
		}
	}
	return fmt.Errorf("binary %s not found in archive", binaryName)
}

func extractZip(archivePath, destDir, binaryName string) error {
	r, err := zip.OpenReader(archivePath)
	if err != nil {
		return err
	}
	defer func() { _ = r.Close() }()

	for _, f := range r.File {
		if filepath.Base(f.Name) == binaryName {
			rc, err := f.Open()
			if err != nil {
				return err
			}
			out, err := os.Create(filepath.Join(destDir, binaryName))
			if err != nil {
				_ = rc.Close()
				return err
			}
			if _, err := io.Copy(out, rc); err != nil { //nolint:gosec
				_ = out.Close()
				_ = rc.Close()
				return err
			}
			_ = out.Close()
			_ = rc.Close()
			return nil
		}
	}
	return fmt.Errorf("binary %s not found in archive", binaryName)
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() { _ = in.Close() }()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() { _ = out.Close() }()

	_, err = io.Copy(out, in)
	return err
}
