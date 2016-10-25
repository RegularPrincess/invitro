// invitro project main.go
package invitro_parser

func main() {
	link := "/analizes/for-doctors/"
	dbConnInf := "postgres://rsulhgnyrdtrhw:AvNZ5aCAKzBbsQGAD8g1er3Ikd@ec2-54-75-228-77.eu-west-1.compute.amazonaws.com:5432/ddusrnru9159g2"
	Scrape(link, dbConnInf)
}
