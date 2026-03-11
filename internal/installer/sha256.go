package installer

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func verifySHA256(checksumsPath, binPath, assetName string) error {
	f, err := os.Open(checksumsPath)
	if err != nil {
		return err
	}
	defer f.Close()

	sumWant := ""
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}
		name := parts[len(parts)-1]
		name = strings.TrimPrefix(name, "*")
		name = filepath.Base(name)
		if name == assetName {
			sumWant = parts[0]
			break
		}
	}
	if sumWant == "" {
		return fmt.Errorf("checksum entry not found for %s", assetName)
	}

	b, err := os.ReadFile(binPath)
	if err != nil {
		return err
	}
	sum := sha256.Sum256(b)
	sumGot := hex.EncodeToString(sum[:])
	if !strings.EqualFold(sumGot, sumWant) {
		return fmt.Errorf("checksum mismatch: got %s want %s", sumGot, sumWant)
	}
	return nil
}
