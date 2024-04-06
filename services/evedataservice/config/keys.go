package config

import "github.com/andrewapj/dotenvconfig"

func ConfigPathKey() string {
	return "APP_CONFIG_PATH"
}

func LogLevel() string {
	return dotenvconfig.GetKeyMust("LOG_LEVEL")
}

func EsiConcurrency() int {
	return dotenvconfig.GetKeyAsIntMust("ESI_CONCURRENCY")
}

func EsiDateAdditionalTime() int {
	return dotenvconfig.GetKeyAsIntMust("ESI_DATE_ADDITIONAL_TIME")
}

func EsiDatasource() string {
	return dotenvconfig.GetKeyMust("ESI_DATASOURCE")
}

func EsiDateLayout() string {
	return dotenvconfig.GetKeyMust("ESI_DATE_LAYOUT")
}

func EsiDomain() string {
	return dotenvconfig.GetKeyMust(EsiDomainKey())
}

func EsiDomainKey() string {
	return "ESI_DOMAIN"
}

func EsiHeaderExpiresKey() string {
	return dotenvconfig.GetKeyMust("ESI_HEADER_EXPIRES_KEY")
}

func EsiHeaderPagesKey() string {
	return dotenvconfig.GetKeyMust("ESI_HEADER_PAGES_KEY")
}

func EsiLanguage() string {
	return dotenvconfig.GetKeyMust("ESI_LANGUAGE")
}

func EsiProtocol() string {
	return dotenvconfig.GetKeyMust(EsiProtocolKey())
}

func EsiProtocolKey() string {
	return "ESI_PROTOCOL"
}

func EsiTimeout() int {
	return dotenvconfig.GetKeyAsIntMust("ESI_TIMEOUT")
}

func EsiUserAgent() string {
	return dotenvconfig.GetKeyMust("ESI_USER_AGENT")
}
