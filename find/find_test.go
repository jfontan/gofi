package find

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFind(t *testing.T) {
	tmp, clean, expected := prepareFindTmp(t, 10, 100, 2)
	defer clean()

	f := New(tmp, 1)
	files, err := f.Find()
	require.NoError(t, err)

	require.Equal(t, expected, cleanFiles(tmp, files))
}

func prepareFindTmp(
	t *testing.T,
	dirs, files, depth int,
) (string, func(), []string) {
	t.Helper()

	path, err := os.MkdirTemp("", "gofind-")
	require.NoError(t, err)

	generated, err := mkFindDir(path, dirs, files, depth)
	require.NoError(t, err)

	return path, func() {
		err := os.RemoveAll(path)
		require.NoError(t, err)
	}, cleanFiles(path, generated)
}

func mkFindDir(path string, dirs, files, depth int) ([]string, error) {
	if depth < 0 {
		return nil, nil
	}

	var generated []string

	for i := 0; i < dirs; i++ {
		p := filepath.Join(path, fmt.Sprintf("d%d", i))
		err := os.MkdirAll(p, 0770)
		if err != nil {
			return nil, err
		}
		f, err := mkFindDir(p, dirs, files, depth-1)
		if err != nil {
			return nil, err
		}
		generated = append(generated, f...)

		p = filepath.Join(path, fmt.Sprintf(".d%d", i))
		err = os.MkdirAll(p, 0770)
		if err != nil {
			return nil, err
		}
		f, err = mkFindDir(p, dirs, files, depth-1)
		if err != nil {
			return nil, err
		}
		generated = append(generated, f...)
	}

	for i := 0; i < files; i++ {
		p := filepath.Join(path, fmt.Sprintf("f%d", i))
		err := mkFile(p)
		if err != nil {
			return nil, err
		}
		generated = append(generated, p)

		p = filepath.Join(path, fmt.Sprintf("f%d.ext", i))
		err = mkFile(p)
		if err != nil {
			return nil, err
		}
		generated = append(generated, p)

		p = filepath.Join(path, fmt.Sprintf(".f%d", i))
		err = mkFile(p)
		if err != nil {
			return nil, err
		}
		generated = append(generated, p)
	}

	return generated, nil
}

func mkFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString("text")
	return err
}

func cleanFiles(tmp string, files []string) []string {
	cleaned := make([]string, len(files))
	for i := range files {
		cleaned[i] = strings.TrimPrefix(files[i], tmp)
	}

	sort.Strings(cleaned)
	return cleaned
}
