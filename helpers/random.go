package helpers

import (
	"bufio"
	"math/rand"
	"os"

	"zntr.io/typogenerator"
	"zntr.io/typogenerator/strategy"
)

var Usernames []string

func LoadUsernames() {
	cFile, err := os.Open("username.txt")
	if err != nil {
	}
	defer cFile.Close()

	scanner := bufio.NewScanner(cFile)

	for scanner.Scan() {
		Usernames = append(Usernames, scanner.Text())
	}

}
func GetRealisticUsername() string {
	var woo []string

	username := Usernames[rand.Intn(len(Usernames))]

	all := []strategy.Strategy{
		strategy.Omission,
		strategy.Repetition,
		strategy.VowelSwap,
		strategy.Addition,
	}

	results, err := typogenerator.Fuzz(username, all...)
	if err != nil {
		return ""
	}

	for _, r := range results {
		woo = append(woo, r.Permutations...)
	}

	return woo[rand.Intn(len(woo))]
}

func GenerateRandomString(length int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyz1234567890ABCDEFGHIJKLMNOPKRSTUVWXYZ")
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func RandomUsername() string {
	username := GenerateRandomString(15)
	return username
}
