package main

import (
	"fmt"

	"github.com/MarkSaravi/drone-go/types"
)

func dataQualityReport(readingData types.ImuReadingQualities) {
	fmt.Println("total data:             ", readingData.Total)
	fmt.Println("number of bad imu data: ", readingData.BadData)
	fmt.Println("number of bad timing:   ", readingData.BadInterval)
	fmt.Println("bad timing rate:        ", float64(readingData.BadInterval)/float64(readingData.Total)*100)
	fmt.Println("max bad timing:         ", readingData.MaxBadInterval)
	fmt.Println("Program stopped.")

}
