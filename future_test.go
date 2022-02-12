package future


import (
	"testing"
  "fmt"
  "time"
)

func TestSum(t *testing.T) {
  fut1 := InstantFutureOf(15)
  fmt.Println(fut1.GetNow())
  fut1.Then(func(i int) {
		fmt.Printf("fut1 is %d\n", i)
	}).Then(func (i int) {
		fmt.Printf("I say again, fut1 is ready and value is %d!\n", i)
	})

  future2 := DelayedFutureOf("thank you", 5 * time.Second)
  future2.Then(func(s string){fmt.Println("You said", s)})
  for {
    ok, v := future2.GetTimeout(300*time.Millisecond)
    if ok {
			fmt.Println("fut2 is ready:", v)
      break
		} else {
			fmt.Printf("Not ready yet:fut2\n")
			fmt.Printf("Launched/Active/Exit %d %d %d\n", LaunchCount(), ActiveCount(), ExitCount())
		}
	}
	fmt.Printf("Launched/Active/Exit %d %d %d\n", LaunchCount(), ActiveCount(), ExitCount())


  future3 := DelayedFutureOf([]int{1,2,3,4,5}, 5*time.Second)

	future4 := Chain(future3, func(arr []int) int {
		sum :=0
		for _, v := range(arr) {
			sum +=v
		}
		return sum
	})

	future4.Then(func (i int) {
		fmt.Printf("The array sum is %d\n", i)
	})

	fmt.Println(future4.GetWait())
}
