package cli

import (
	"fmt"
	"bufio"
	"os"
	"strconv"
)

func Unswer(max int, caption string) int {
	fmt.Print(caption)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	text := scanner.Text()
	if text == "" {
		return 0
	}
	selected, err := strconv.Atoi(text)
	if err != nil || selected > max || selected <= 0 {
		return MustUnswer(max, "Wrong input, try again: ")
	}
	return selected
}
func MustUnswer(max int, caption string) int {
	selected := Unswer(max, caption)
	if selected == 0 {
		return MustUnswer(max, caption)
	}
	return selected
}
func StringUnswer(caption string) string {
	fmt.Print(caption)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}
