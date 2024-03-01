package calc

import (
	"fmt"
	"log"
	"testing"
)

func checkInfo(t *testing.T, infos []iterInfo, i int, expectedDim int, expectedSkip uint64, expectedLimit uint64) {
	info := infos[i]
	t.Logf(fmt.Sprintf("%v: %s", i, info.String()))
	if info.dim != expectedDim {
		t.Fatalf("bad dim in info[%v]: expected %v, got %v", i, expectedDim, info.dim)
	}
	if info.needSkip {
		if info.skip != expectedSkip {
			t.Fatalf("bad skip length in info[%v]: expected %v, got %v", i, expectedSkip, info.skip)
		}
	}
	if info.needLimit {
		if info.limit != expectedLimit {
			t.Fatalf("bad limit in info[%v]: expected %v, got %v", i, expectedLimit, info.limit)
		}
	}
}

func TestPrivate_calcIterationInfos(t *testing.T) {
	var alphabetLength uint64 = 36
	var startIndex uint64 = 0
	var partCount uint64 = 1000
	maxLength := 2

	infos := calcIterationInfos(log.Default(), alphabetLength, maxLength, startIndex, partCount)
	if len(infos) != 2 {
		t.Fatalf("expected 2 iter infos, got %v", infos)
	}

	checkInfo(t, infos, 0, 1, 0, 0)
	checkInfo(t, infos, 1, 2, 0, 1000-36)

	startIndex = 1
	infos = calcIterationInfos(log.Default(), alphabetLength, maxLength, startIndex, partCount)
	if len(infos) != 2 {
		t.Fatalf("expected 2 iter infos, got %v", infos)
	}

	checkInfo(t, infos, 0, 1, 1, 35)
	checkInfo(t, infos, 1, 2, 0, 1000-35)
}
