package pqueue

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/linxGnu/pqueue/common"
	"github.com/stretchr/testify/require"
)

func TestLoadInfos(t *testing.T) {
	t.Run("Invalid", func(t *testing.T) {
		_, err := loadFileInfos("/abc", nil)
		require.Error(t, err)
	})

	t.Run("Happy", func(t *testing.T) {
		f1, err := ioutil.TempFile(tmpDir, "seg_")
		require.NoError(t, err)
		require.NoError(t, f1.Close())

		f2, err := ioutil.TempFile(tmpDir, "seg_")
		require.NoError(t, err)
		require.NoError(t, f2.Close())

		_, err = loadFileInfos(tmpDir, func(os.DirEntry) (os.FileInfo, error) {
			return nil, fmt.Errorf("fake error")
		})
		require.Error(t, err)

		files, err := loadFileInfos(tmpDir, fileInfoExtractor)
		require.NoError(t, err)
		for i := range files {
			_ = os.Remove(files[i].path)
		}
	})
}

func TestLoading(t *testing.T) {
	t.Run("Error", func(t *testing.T) {
		_, err := load(QueueSettings{
			DataDir: "/abc",
		}, nil)
		require.Error(t, err)

		_, err = load(QueueSettings{
			DataDir:     tmpDir,
			EntryFormat: 123,
		}, &segmentHeader{})
		require.Error(t, err)

		_, err = load(QueueSettings{
			DataDir:       tmpDir,
			SegmentFormat: 123,
		}, &segmentHeader{})
		require.Error(t, err)
	})

	t.Run("OK", func(t *testing.T) {
		q, err := load(QueueSettings{
			DataDir:       tmpDir,
			SegmentFormat: common.SegmentV1,
			EntryFormat:   common.EntryV1,
		}, &segmentHeader{})
		require.NoError(t, err)
		for q.segments.Len() > 0 {
			front := q.segments.Front()
			_ = os.Remove(front.Value.(*segment).path)
			q.segments.Remove(front)
		}
	})
}
