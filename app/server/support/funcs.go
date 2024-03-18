package support

import "log"

func FatalErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func PanicErr(err error) {
	if err != nil {
		panic(err)
	}
}

func WarningErr(err error) {
	if err != nil {
		log.Println("Warning: ", err)
	}
}
