package tool

import (
	"fmt"
	"os"
)

var (
	ffmpegPath string
	ytDlPath   string
	outputPath string
	busCs      string
	dbCs       string
	apiUrl     string
)

const (
	ytdPathEnv       = "YTD_PATH"
	ffmpegPathEnv    = "FFMPEG_PATH"
	bydOutputPathEnv = "BYD_OUTPUT_PATH"
	bydBusEnv        = "BYD_BUS_CS"
	bydDbEnv         = "BYD_DB_CS"
	bydApiEnv        = "BYD_API_URL"
)

func GetFFmpegPath() (string, error) {
	return getValue(ffmpegPath, ffmpegPathEnv)
}

func GetYtdPath() (string, error) {
	return getValue(ytDlPath, ytdPathEnv)
}

func GetOutputPath() (string, error) {
	return getValue(outputPath, bydOutputPathEnv)
}

func GetBusConnectionString() (string, error) {
	return getValue(busCs, bydBusEnv)
}

func GetDbConnectionString() (string, error) {
	return getValue(dbCs, bydDbEnv)
}

func GetApiUrl() (string, error) {
	return getValue(apiUrl, bydApiEnv)
}

func getValue(currentValue string, env string) (string, error) {
	if currentValue != "" {
		return currentValue, nil
	}
	var err error
	currentValue, err = getEnvVar(env)
	if err != nil {
		return "", err
	}
	return currentValue, nil
}

func getEnvVar(env string) (string, error) {
	path, ok := os.LookupEnv(env)
	if !ok {
		return "", fmt.Errorf("%s env var not found", env)
	}
	return path, nil
}
