package config //설정파일(toml) 읽기, 설정값(데이터베이스 접속정보, 카프카 서버 정보) 저장, 중요한 정보 저장 관리

import (
	"github.com/naoina/toml"  // TOML 파일을 파싱하기 위한 라이브러리
    "os"                      // 파일 시스템 작업을 위한 패키지
)

type Config struct {
    DB struct {        // DB 관련 설정을 담는 중첩 구조체
        Database string  // 데이터베이스 이름
        URL      string  // 데이터베이스 연결 URL
    }

    Kafka struct {     // Kafka 관련 설정을 담는 중첩 구조체
        URL      string  // Kafka 서버 URL
        ClientID string  // Kafka 클라이언트 ID
    }
}
func NewConfig(path string) *Config {
    c := new(Config)    // 새로운 Config 구조체 포인터 생성

    if f, err := os.Open(path); err != nil {    // 설정 파일 열기
        panic(err)    // 파일 열기 실패시 패닉
    } else if err = toml.NewDecoder(f).Decode(c); err != nil {    // TOML 파일 디코딩
        panic(err)    // 디코딩 실패시 패닉
    } else {
        return c      // 성공시 설정 반환
    }
}