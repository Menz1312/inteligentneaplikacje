package db

func init() {
	err := setup("root:@tcp(127.0.0.1:3306)/")
	if err != nil {
		panic(err)
	}
	checkDatabase("iai")
}
