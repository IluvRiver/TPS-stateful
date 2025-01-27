package config

import (
	"os"

	"github.com/naoina/toml"
)

type Config struct {
	DB struct { // DB 관련 설정을 담는 중첩 구조체
		Database string // 데이터베이스 이름
		URL      string // 데이터베이스 연결 URL
	}

	Kafka struct {
		URL     string
		GroupID string
	}
	Info struct { //포트관리
		Port string
	}
}

func NewConfig(path string) *Config {
	c := new(Config) // 새로운 Config 구조체 포인터 생성

	if f, err := os.Open(path); err != nil { // 설정 파일 열기
		panic(err) // 파일 열기 실패시 패닉
	} else if err = toml.NewDecoder(f).Decode(c); err != nil { // TOML 파일 디코딩
		panic(err) // 디코딩 실패시 패닉
	} else {
		return c // 성공시 설정 반환
	}
}
