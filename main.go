package main

import (
	"fmt"
	"time"

	"github.com/asticode/go-astisub"
)

func main() {
	// Open subtitles
	s1, err := astisub.OpenFile("srtSubtitles.srt")
	if err != nil {
		fmt.Println(err.Error())
	}

	// Wrap subtitles
	exS1 := &extendedSubtitles{s1}

	// Fragment the subtitles
	fragmentedItems := exS1.ExFragment(10 * time.Minute)

	// fragmentedItemsをloopさせる
	for i, items := range fragmentedItems {
		astisub.Subtitles{Items: items}.Write(fmt.Sprintf("new_srtSubtitles_%d.srt", i))
	}
}

type extendedSubtitles struct {
	*astisub.Subtitles
}

// 参照（以下URLのFragmentメソッドを拡張したもの）
// https://github.com/asticode/go-astisub/blob/721d3fc8258bcc5912423c551da15ad39ebcc553/subtitles.go#L569
func (es *extendedSubtitles) ExFragment(f time.Duration) [][]*astisub.Item {
	// Nothing to fragment
	if len(es.Items) == 0 {
		return [][]*astisub.Item{}
	}

	// 分割された配列を定義
	var flangments [][]*astisub.Item
	flangmentStartIndex := 0
	// Here we want to simulate fragments of duration f until there are no subtitles left in that period of time
	var fragmentStartAt, fragmentEndAt = time.Duration(0), f
	for fragmentStartAt < es.Items[len(es.Items)-1].EndAt {
		// We loop through subtitles and process the ones that either contain the fragment start at,
		// or contain the fragment end at
		//
		// It's useless processing subtitles contained between fragment start at and end at
		//             |____________________|             <- subtitle
		//           |                        |
		//   fragment start at        fragment end at
		for i, sub := range es.Items {
			// Init
			var newSub = &astisub.Item{}
			*newSub = *sub

			// A switch is more readable here
			switch {
			// Subtitle contains fragment start at
			// |____________________|                         <- subtitle
			//           |                        |
			//   fragment start at        fragment end at
			case sub.StartAt < fragmentStartAt && sub.EndAt > fragmentStartAt:
				sub.StartAt = fragmentStartAt
				newSub.EndAt = fragmentStartAt
			// Subtitle contains fragment end at
			//                         |____________________| <- subtitle
			//           |                        |
			//   fragment start at        fragment end at
			case sub.StartAt < fragmentEndAt && sub.EndAt > fragmentEndAt:
				sub.StartAt = fragmentEndAt
				newSub.EndAt = fragmentEndAt
			case sub.StartAt >= fragmentEndAt && i > 0 && es.Items[i-1].EndAt < fragmentEndAt:
				fmt.Printf("flangment: %v ~ %v\n", flangmentStartIndex, i)
				flangments = append(flangments, es.Items[flangmentStartIndex:i])
				flangmentStartIndex = i
				continue
			default:
				// 既存の値の更新が不要な場合はスキップ
				continue
			}

			// case文に一致して、時間の区切り用に新しいサブタイトルが必要な場合Insert
			es.Items = append(es.Items[:i], append([]*astisub.Item{newSub}, es.Items[i:]...)...)

			// Insert後追加した要素も含めて配列を分ける
			fmt.Printf("flangment: %v ~ %v\n", flangmentStartIndex, i+1)
			flangments = append(flangments, es.Items[flangmentStartIndex:i+1])
			flangmentStartIndex = i + 1
		}

		// Update fragments boundaries
		fragmentStartAt += f
		fragmentEndAt += f
	}

	// Order
	es.Order()

	// 最後のフラグメントを追加
	flangments = append(flangments, es.Items[flangmentStartIndex:])
	return flangments
}
