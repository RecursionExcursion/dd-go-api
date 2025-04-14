package cfbr

import (
	"log"
)

// First accept args (stat weights)/(year)
func CFBR() {

	//Use compressed data or fetch new

	//use old

	//fetch new
	season, err := getNewData("fbs", 2024)
	if err != nil {
		panic(err)
	}

	/* TODO test logs, -rm */

	s, err := season.FindSchoolById(194)
	if err != nil {
		panic(err)
	}
	log.Println(s)
	log.Println(len(s.Games))

}
