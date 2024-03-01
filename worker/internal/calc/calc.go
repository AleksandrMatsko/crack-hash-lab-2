package calc

import (
	"context"
	"crypto/md5"
	"distributed.systems.labs/shared/pkg/alphabet"
	"distributed.systems.labs/shared/pkg/cartesian-gen"
	"distributed.systems.labs/shared/pkg/contracts"
	"fmt"
	"log"
)

type iterInfo struct {
	dim int

	needSkip bool
	skip     uint64

	needLimit bool
	limit     uint64
}

func (info iterInfo) String() string {
	return fmt.Sprintf("dim = %v\n\tneedSkip = %v, skip = %v\n\tneedLimit = %v, limit %v",
		info.dim, info.needSkip, info.skip, info.needLimit, info.limit)
}

func calcIterationInfos(logger *log.Logger, alphabetLength uint64, maxLength int, startIndex uint64, partCount uint64) []iterInfo {
	maxVal := startIndex + partCount
	infos := make([]iterInfo, 0)

	degs := make([]uint64, maxLength)
	sums := make([]uint64, maxLength+1)
	if maxLength <= 0 {
		return infos
	}
	degs[0] = alphabetLength
	sums[0] = 0
	sums[1] = alphabetLength
	for i := 1; i < len(degs); i++ {
		degs[i] = degs[i-1] * alphabetLength
		sums[i+1] = sums[i] + degs[i]
	}
	logger.Printf("degs: %v", degs)
	logger.Printf("sums: %v", sums)

	startDim := 0
	for sums[startDim] <= startIndex {
		startDim += 1
	}
	endDim := startDim
	for sums[endDim] < maxVal {
		endDim += 1
	}
	logger.Printf("startDim = %v endDim = %v", startDim, endDim)
	for i := startDim - 1; i < endDim; i++ {
		info := iterInfo{
			dim: i + 1,
		}
		if startIndex >= sums[i] {
			info.needSkip = true
			info.skip = startIndex - sums[i]
			if sums[i+1] >= maxVal {
				info.needLimit = true
				info.limit = partCount
			} else {
				info.needLimit = false
			}
		} else {
			info.needSkip = false
			if sums[i+1] >= maxVal {
				info.needLimit = true
				info.limit = maxVal - sums[i]
			} else {
				info.needLimit = false
			}
		}
		infos = append(infos, info)
	}
	return infos
}

// ProcessRequest should be called in separate goroutine from endpoint handler
func ProcessRequest(ctx context.Context, req contracts.TaskRequest, resChan chan<- contracts.TaskResultRequest) {
	defaultLogger := log.Default()
	logger := log.New(
		defaultLogger.Writer(),
		fmt.Sprintf("request-id: %s ", req.RequestID),
		defaultLogger.Flags()|log.Lmsgprefix)

	logger.Printf("alphabet = %s startIndex = %v partCount = %v", req.Alphabet, req.StartIndex, req.PartCount)
	A := alphabet.InitAlphabet([]rune(req.Alphabet))
	infos := calcIterationInfos(logger, uint64(A.Length()), req.MaxLength, req.StartIndex, req.PartCount)
	var completed uint64
	cracks := make([]string, 0)
	for _, info := range infos {
		dims := make([]uint64, 0)
		for j := 0; j < info.dim; j++ {
			dims = append(dims, uint64(A.Length()))
		}
		gen := cartesian_gen.NewCartesianGenerator(dims)
		if info.needSkip {
			gen = gen.Skip(info.skip)
		}
		if info.needLimit {
			gen = gen.Limit(info.limit)
		}
		for gen.HasNext() {
			select {
			case <-ctx.Done():
				logger.Printf("ctx.Done")
				return
			default:
				ids := gen.Product()
				word := A.GetWord(ids)
				gotHash := fmt.Sprintf("%x", md5.Sum([]byte(word)))
				if gotHash == req.ToCrack {
					cracks = append(cracks, word)
				}
				completed += 1
				if completed%5_000_000 == 0 {
					logger.Printf("completed %v / %v (start_index = %v)",
						completed, req.PartCount, req.StartIndex)
				}
			}
		}
	}
	logger.Printf("calc finished completed %v / %v, sending results to notifier ...", completed, req.PartCount)
	res := contracts.TaskResultRequest{
		StartIndex: req.StartIndex,
		RequestID:  req.RequestID,
		Cracks:     cracks,
	}
	resChan <- res
}
