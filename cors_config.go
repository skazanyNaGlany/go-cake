package go_cake

import (
	"regexp"
)

type CORSConfig struct {
	GetAllowedOrigins     []string
	DeleteAllowedOrigins  []string
	InsertAllowedOrigins  []string
	UpdateAllowedOrigins  []string
	getCompiledOrigins    []*regexp.Regexp
	deleteCompiledOrigins []*regexp.Regexp
	insertCompiledOrigins []*regexp.Regexp
	updateCompiledOrigins []*regexp.Regexp
}

func NewCORSConfig(
	GetAllowedOrigins []string,
	DeleteAllowedOrigins []string,
	InsertAllowedOrigins []string,
	UpdateAllowedOrigins []string) (*CORSConfig, error) {
	corsConfig := CORSConfig{
		GetAllowedOrigins:    GetAllowedOrigins,
		DeleteAllowedOrigins: DeleteAllowedOrigins,
		InsertAllowedOrigins: InsertAllowedOrigins,
		UpdateAllowedOrigins: UpdateAllowedOrigins}

	for _, pattern := range corsConfig.GetAllowedOrigins {
		compiled, err := regexp.Compile(pattern)

		if err != nil {
			return nil, err
		}

		corsConfig.getCompiledOrigins = append(corsConfig.getCompiledOrigins, compiled)
	}

	for _, pattern := range corsConfig.DeleteAllowedOrigins {
		compiled, err := regexp.Compile(pattern)

		if err != nil {
			return nil, err
		}

		corsConfig.deleteCompiledOrigins = append(corsConfig.deleteCompiledOrigins, compiled)
	}

	for _, pattern := range corsConfig.InsertAllowedOrigins {
		compiled, err := regexp.Compile(pattern)

		if err != nil {
			return nil, err
		}

		corsConfig.insertCompiledOrigins = append(corsConfig.insertCompiledOrigins, compiled)
	}

	for _, pattern := range corsConfig.UpdateAllowedOrigins {
		compiled, err := regexp.Compile(pattern)

		if err != nil {
			return nil, err
		}

		corsConfig.updateCompiledOrigins = append(corsConfig.updateCompiledOrigins, compiled)
	}

	return &corsConfig, nil
}

func NewDefaultCORSConfig() (*CORSConfig, error) {
	return NewCORSConfig(
		[]string{".*"},
		[]string{".*"},
		[]string{".*"},
		[]string{".*"})
}
