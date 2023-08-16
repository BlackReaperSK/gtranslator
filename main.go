package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/bregydoc/gtranslate"
)

func main() {
	var toLang string

	// Definir flags para receber o valor do idioma alvo
	flag.StringVar(&toLang, "l", "en", "Código de idioma alvo (ex: en, ja, es, fr, etc.)")
	flag.Parse()

	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		var wg sync.WaitGroup
		semaphore := make(chan struct{}, 5) // Limitar a 5 goroutines
		scanner := bufio.NewScanner(os.Stdin)

		for scanner.Scan() {
			text := scanner.Text()
			wg.Add(1)
			semaphore <- struct{}{}
			go func(text string) {
				defer func() {
					<-semaphore
					wg.Done()
				}()
				translated, err := gtranslate.TranslateWithParams(
					text,
					gtranslate.TranslationParams{
						From: "auto",
						To:   toLang,
					},
				)
				if err != nil {
					panic(err)
				}
				fmt.Printf("%s \n", translated)

				// Adicionar atraso entre as chamadas
				time.Sleep(120 * time.Millisecond)
			}(text)
		}

		wg.Wait()
	}
}
