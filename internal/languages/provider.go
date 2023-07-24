package languages

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"golang.org/x/exp/slog"
)

type LanguageProvider interface {
	Sync() error
	GetLanguage(id string) (ProgrammingLanguage, error)
	GetLanguages() ([]ProgrammingLanguage, error)
	FindByFileExtension(extension string) (ProgrammingLanguage, error)
}

type LanguageProviderBase struct {
	languages []ProgrammingLanguage
}

type JsonLanguageProvider struct {
	LanguageProviderBase
	jsonPath string
}

func NewJsonLanguageProvider(jsonPath string) (*JsonLanguageProvider, error) {
    var result = new(JsonLanguageProvider)
    result.jsonPath = jsonPath
    err := result.Sync()
    if err != nil {
        return nil, err
    }
    return result, nil
}

func (provider *JsonLanguageProvider) Sync() error {
	jsonDataBytes, err := os.ReadFile(provider.jsonPath)
	if err != nil {
		return err
	}
	
    var languages []ProgrammingLanguage
    err = json.Unmarshal(jsonDataBytes, &languages)
    if err != nil {
        return err
    }
    
    provider.languages = languages
    return nil
}

func (base *LanguageProviderBase) GetLanguage(id string) (ProgrammingLanguage, error) {
	for _, language := range base.languages {
		if language.Id == id {
			return language, nil
		}
	}
	return ProgrammingLanguage{}, errors.New("language not found")
}

func (base *LanguageProviderBase) GetLanguages() ([]ProgrammingLanguage, error) {
	return base.languages, nil
}

func (base *LanguageProviderBase) FindByFileExtension(extension string) (ProgrammingLanguage, error) {
    slog.Info("finding language by extension", slog.String("extension", extension))
	var results []ProgrammingLanguage
	for _, language := range base.languages {
		langExtension := filepath.Ext(language.CodeFilename)
		if langExtension == extension {
			results = append(results, language)
		}
	}
    slog.Info("found languages", slog.Int("count", len(results)))
	if len(results) == 1 {
		return results[0], nil
	}
	if len(results) > 1 {
		return ProgrammingLanguage{}, errors.New("multiple languages found")
	}
	if len(results) == 0 {
		return ProgrammingLanguage{}, errors.New("language not found")
	}
	return ProgrammingLanguage{}, errors.New("unknown error")
}
