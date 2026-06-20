package config

import (
	"os"
)

// Config รวม env ทั้งหมดของแอป โหลดครั้งเดียวตอน startup
type Config struct {
	Port            string
	MongoURI        string
	MongoDB         string
	JWTSecret       string
	OpenRouterKey   string
	OpenRouterModel string
}

// Load อ่านค่าจาก environment variables พร้อม default ที่ใช้ตอน dev
func Load() Config {
	return Config{
		Port:            getEnv("PORT", "8080"),
		MongoURI:        getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:         getEnv("MONGO_DB", "kaikaa"),
		JWTSecret:       getEnv("JWT_SECRET", "dev-secret-change-me"),
		OpenRouterKey:   getEnv("OPENROUTER_KEY", ""),
		OpenRouterModel: getEnv("OPENROUTER_MODEL", "openai/gpt-4o-mini"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
