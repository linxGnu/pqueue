package pqueue

import (
	"container/list"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type file struct {
	modTime time.Time
	path    string
}

func load(settings QueueSettings, segHeader segmentHeadWriter) (*queue, error) {
	if settings.MaxEntriesPerSegment <= 0 {
		settings.MaxEntriesPerSegment = DefaultMaxEntriesPerSegment
	}

	files, err := loadFileInfos(settings.DataDir, fileInfoExtractor)
	if err != nil {
		return nil, err
	}

	segments := list.New()
	for i := range files {
		segments.PushBack(&segment{
			readable: false,
			path:     files[i].path,
		})
	}

	// create new segment for upcoming entries
	q := &queue{
		settings:      settings,
		segHeadWriter: segHeader,
	}

	seg, err := q.newSegment()
	if err != nil {
		return nil, err
	}
	segments.PushBack(seg)

	q.segments = segments
	return q, nil
}

func loadFileInfos(dir string, infoExtractor func(os.DirEntry) (os.FileInfo, error)) ([]file, error) {
	fileList, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	files := make([]file, 0, len(fileList))
	for i := range fileList {
		if strings.HasPrefix(fileList[i].Name(), segPrefix) {
			info, e := infoExtractor(fileList[i])
			if e != nil {
				return nil, e
			}

			files = append(files, file{
				path:    filepath.Join(dir, fileList[i].Name()),
				modTime: info.ModTime(),
			})
		}
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].modTime.Before(files[j].modTime)
	})

	return files, nil
}

func fileInfoExtractor(f os.DirEntry) (os.FileInfo, error) {
	return f.Info()
}
