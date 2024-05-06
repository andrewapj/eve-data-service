package config

import "github.com/andrewapj/dotenvconfig"

func ConfigPathKey() string { return "APP_CONFIG_PATH" }

func DbUrl() string { return dotenvconfig.GetKeyMust(DbUrlKey()) }

func DbUrlKey() string { return "DB_URL" }

func EsiConcurrency() int { return dotenvconfig.GetKeyAsIntMust("ESI_CONCURRENCY") }

func EsiDateDefaultAdditionalTimeSeconds() int {
	return dotenvconfig.GetKeyAsIntMust("ESI_DATE_DEFAULT_ADDITIONAL_TIME_SECONDS")
}

func EsiDatasource() string { return dotenvconfig.GetKeyMust("ESI_DATASOURCE") }

func EsiDateLayout() string { return dotenvconfig.GetKeyMust("ESI_DATE_LAYOUT") }

func EsiDomain() string { return dotenvconfig.GetKeyMust(EsiDomainKey()) }

func EsiDomainKey() string { return "ESI_DOMAIN" }

func EsiHeaderExpiresKey() string { return dotenvconfig.GetKeyMust("ESI_HEADER_EXPIRES_KEY") }

func EsiHeaderPagesKey() string { return dotenvconfig.GetKeyMust("ESI_HEADER_PAGES_KEY") }

func EsiHeaderUserAgentKey() string { return dotenvconfig.GetKeyMust("ESI_HEADER_USER_AGENT") }

func EsiLanguage() string { return dotenvconfig.GetKeyMust("ESI_LANGUAGE") }

func EsiMaxRetries() int { return dotenvconfig.GetKeyAsIntMust("ESI_MAX_RETRIES") }

func EsiProtocol() string { return dotenvconfig.GetKeyMust(EsiProtocolKey()) }

func EsiProtocolKey() string { return "ESI_PROTOCOL" }

func EsiTimeout() int { return dotenvconfig.GetKeyAsIntMust("ESI_TIMEOUT") }

func EsiUserAgent() string { return dotenvconfig.GetKeyMust("ESI_USER_AGENT") }

func LoaderIntervalSeconds() int { return dotenvconfig.GetKeyAsIntMust("LOADER_INTERVAL_SECONDS") }

func LogLevel() string { return dotenvconfig.GetKeyMust("LOG_LEVEL") }
