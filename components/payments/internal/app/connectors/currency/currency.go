package currency

import (
	"fmt"
	"strings"

	"github.com/formancehq/payments/internal/app/models"
)

type currency struct {
	decimals int
}

func currencies() map[string]currency {
	return map[string]currency{
		"AED": {2}, //  UAE Dirham
		"AMD": {2}, //  Armenian Dram
		"ANG": {2}, //  Netherlands Antillian Guilder
		"ARS": {2}, //  Argentine Peso
		"ATS": {2}, //  Euro
		"AUD": {2}, //  Australian Dollar
		"AWG": {2}, //  Aruban Guilder
		"BAM": {2}, //  Bosnia and Herzegovina, Convertible Marks
		"BDT": {2}, //  Bangladesh, Taka
		"BEF": {2}, //  Euro
		"BHD": {3}, //  Bahraini Dinar
		"BMD": {2}, //  Bermudian Dollar
		"BND": {2}, //  Brunei Dollar
		"BOB": {2}, //  Bolivia, Boliviano
		"BRL": {2}, //  Brazilian Real
		"BSD": {2}, //  Bahamian Dollar
		"BWP": {2}, //  Botswana, Pula
		"BZD": {2}, //  Belize Dollar
		"CAD": {2}, //  Canadian Dollar
		"CHF": {2}, //  Swiss Franc
		"CLP": {0}, //  Chilean Peso
		"CNY": {2}, //  China Yuan Renminbi
		"COP": {2}, //  Colombian Peso
		"CRC": {2}, //  Costa Rican Colon
		"CUC": {2}, //  Cuban Convertible Peso
		"CUP": {2}, //  Cuban Peso
		"CYP": {2}, //  Cyprus Pound
		"CZK": {2}, //  Czech Koruna
		"DEM": {2}, //  Euro
		"DKK": {2}, //  Danish Krone
		"DOP": {2}, //  Dominican Peso
		"EEK": {2}, //  Euro
		"EGP": {2}, //  Egyptian Pound
		"ESP": {2}, //  Euro
		"EUR": {2}, //  Euro
		"FIM": {2}, //  Euro
		"FRF": {2}, //  Euro
		"GBP": {2}, //  Pound Sterling
		"GHC": {2}, //  Ghana, Cedi
		"GIP": {2}, //  Gibraltar Pound
		"GRD": {2}, //  Euro
		"GTQ": {2}, //  Guatemala, Quetzal
		"HKD": {2}, //  Hong Kong Dollar
		"HNL": {2}, //  Honduras, Lempira
		"HRK": {2}, //  Croatian Kuna
		"HUF": {0}, //  Hungary, Forint
		"IDR": {2}, //  Indonesia, Rupiah
		"IEP": {2}, //  Euro
		"ILS": {2}, //  New Israeli Shekel
		"INR": {2}, //  Indian Rupee
		"IRR": {2}, //  Iranian Rial
		"ISK": {0}, //  Iceland Krona
		"ITL": {2}, //  Euro
		"JMD": {2}, //  Jamaican Dollar
		"JOD": {3}, //  Jordanian Dinar
		"JPY": {0}, //  Japan, Yen
		"KES": {2}, //  Kenyan Shilling
		"KRW": {0}, //  South Korea, Won
		"KWD": {3}, //  Kuwaiti Dinar
		"KYD": {2}, //  Cayman Islands Dollar
		"LBP": {0}, //  Lebanese Pound
		"LTL": {2}, //  Lithuanian Litas
		"LUF": {2}, //  Euro
		"LVL": {2}, //  Latvian Lats
		"MKD": {2}, //  Macedonia, Denar
		"MTL": {2}, //  Maltese Lira
		"MUR": {0}, //  Mauritius Rupee
		"MXN": {2}, //  Mexican Peso
		"MYR": {2}, //  Malaysian Ringgit
		"MZM": {2}, //  Mozambique Metical
		"NLG": {2}, //  Euro
		"NOK": {2}, //  Norwegian Krone
		"NPR": {2}, //  Nepalese Rupee
		"NZD": {2}, //  New Zealand Dollar
		"OMR": {3}, //  Rial Omani
		"PEN": {2}, //  Peru, Nuevo Sol
		"PHP": {2}, //  Philippine Peso
		"PKR": {2}, //  Pakistan Rupee
		"PLN": {2}, //  Poland, Zloty
		"PTE": {2}, //  Euro
		"ROL": {2}, //  Romania, Old Leu
		"RON": {2}, //  Romania, New Leu
		"RUB": {2}, //  Russian Ruble
		"SAR": {2}, //  Saudi Riyal
		"SEK": {2}, //  Swedish Krona
		"SGD": {2}, //  Singapore Dollar
		"SIT": {2}, //  Slovenia, Tolar
		"SKK": {2}, //  Slovak Koruna
		"SVC": {2}, //  El Salvador Colon
		"SZL": {2}, //  Swaziland, Lilangeni
		"THB": {2}, //  Thailand, Baht
		"TOP": {2}, //  Tonga, Paanga
		"TRY": {2}, //  New Turkish Lira
		"TZS": {2}, //  Tanzanian Shilling
		"UAH": {2}, //  Ukraine, Hryvnia
		"USD": {2}, //  US Dollar
		"UYU": {2}, //  Peso Uruguayo
		"VEB": {2}, //  Venezuela, Bolivar
		"VEF": {2}, //  Venezuela Bolivares Fuertes
		"VND": {0}, //  Viet Nam, Dong
		"VUV": {0}, //  Vanuatu, Vatu
		"XCD": {2}, //  East Caribbean Dollar
		"ZAR": {2}, //  South Africa, Rand
		"ZWD": {2}, //  Zimbabwe Dollar
	}
}

func FormatAsset(cur string) models.Asset {
	asset := strings.ToUpper(string(cur))

	def, ok := currencies()[asset]
	if !ok {
		return models.Asset(asset)
	}

	if def.decimals == 0 {
		return models.Asset(asset)
	}

	return models.Asset(fmt.Sprintf("%s/%d", asset, def.decimals))
}

func GetPrecision(cur string) int {
	asset := strings.ToUpper(string(cur))

	def, ok := currencies()[asset]
	if !ok {
		return 0
	}

	return def.decimals
}
