package util

func Init() {
	err := initIdGen()
	if err != nil {
		panic(err)
	}
}
