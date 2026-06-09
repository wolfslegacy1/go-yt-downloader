package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/kkdai/youtube/v2"
)

func main() {

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("🔗 Cole a URL do vídeo do YouTube e aperte Enter: ")


	scanner.Scan()
	videoURL := strings.TrimSpace(scanner.Text())

	if videoURL == "" {
		fmt.Println("❌ Erro: Você precisa colar uma URL válida.")
		return
	}

	outputFolder := "downloads"
	err := os.MkdirAll(outputFolder, os.ModePerm)
	if err != nil {
		fmt.Printf("❌ Erro ao criar a pasta '%s': %v\n", outputFolder, err)
		return
	}

	fmt.Println("🕵️ Buscando informações do vídeo...")
	client := youtube.Client{}

	video, err := client.GetVideo(videoURL)
	if err != nil {
		fmt.Printf("❌ Erro ao buscar vídeo: %v\n", err)
		return
	}

	
	formats := video.Formats.WithAudioChannels()
	if len(formats) == 0 {
		fmt.Println("❌ Nenhum formato compatível encontrado.")
		return
	}
	format := &formats[0]

	fmt.Printf("📥 Baixando: '%s'\n", video.Title)
	fmt.Printf("Qualidade: %s | Formato: %s\n", format.QualityLabel, format.MimeType)


	stream, _, err := client.GetStream(video, format)
	if err != nil {
		fmt.Printf("❌ Erro ao obter stream: %v\n", err)
		return
	}
	defer stream.Close()


	filename := fmt.Sprintf("%s.mp4", video.Title)
	finalPath := filepath.Join(outputFolder, filename)

	file, err := os.Create(finalPath)
	if err != nil {
		fmt.Printf("❌ Erro ao criar arquivo: %v\n", err)
		return
	}
	defer file.Close()


	_, err = io.Copy(file, stream)
	if err != nil {
		fmt.Printf("❌ Erro ao salvar o vídeo: %v\n", err)
		return
	}

	fmt.Printf("🎉 Download concluído com sucesso! Salvo em: %s\n", finalPath)
}
