package main

import (
	"io/ioutil"
	"fmt"
	"strings"
	"unicode"
	"time"
	"log"
	"sync"
)

func main() {
	dictionary := string(loadFile("dictionaries/large"))
	dictionaryMap := makeDictionaryMap(dictionary)
	text := strings.ToLower(string(loadFile("texts/austinpowers.txt")))
	textArr:= getWordsFromText(text)

	fmt.Println("misspellings: ",getMisspellings(dictionaryMap,textArr))

}
//================================================================================================
func loadFile(filename string) []byte{
	defer timeTrack(time.Now(), "load "+filename)
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Print(err)
	}

	return file
}
//================================================================================================
func makeDictionaryMap(text string) map[string]string {
	defer timeTrack(time.Now(), "make dictionary")
	wordsMap := make(map[string]string)
	s := strings.Split(text, "\n")
	for i:= range s{
		wordsMap[s[i]] = s[i]
	}
	return wordsMap
}
//================================================================================================
func getWordsFromText(text string) []string {
	var words []string
	word := ""
	for _, value := range text {
		if unicode.IsLetter(value) || (value == '\'' && len(word) > 0) {
			word += string(value)
		} else if unicode.IsDigit(value) {
			word = ""
		} else if len(word) > 0 {
			words = append(words,word)
			word = ""
		}
	}

	return words
}
//================================================================================================
func getMisspellings(dictionary map[string]string, text []string) int {
	defer timeTrack(time.Now(), "misspellings counting")
	misspellings := 0
	m := sync.Mutex{}
	m.Lock()
	c := sync.NewCond(&m)

	for _,v := range text{
		go func(v string){
			m.Lock()
			defer m.Unlock()
			if dictionary[v] != v{
				misspellings++
			}
			c.Broadcast()
		}(v)
	}
	c.Wait()
	m.Unlock()
	return misspellings
}
//================================================================================================
func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}