package main
import "fmt"

func main(){
	fmt.Println(doubleScore(0))
	fmt.Println(doubleScore(20))
	fmt.Println(doubleScore(100))
}
func doubleScore(source int) (score int) {
	defer func(){
		if score<1 || score >= 100 {
			score = source
		}
	}()
	score = source * 2
	return 
}