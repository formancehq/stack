package currency

import (
	"errors"
	"fmt"
	"strings"
)

var (
	UnsupportedCurrencies = map[string]struct{}{
		"HUF": {},
		"ISK": {},
		"TWD": {},
	}

	ISO4217Currencies = map[string]int{
		"AFN": 2, //  Afghan afghani
		"EUR": 2, //  Euro
		"ALL": 2, //  Albanian lek
		"DZD": 2, //  Algerian dinar
		"USD": 2, //  United States dollar
		"AOA": 2, //  Angolan kwanza
		"XCD": 2, //  East Caribbean dollar
		"ARS": 2, //  Argentine peso
		"AMD": 2, //  Armenian dram
		"AWG": 2, //  Aruban florin
		"AUD": 2, //  Australian dollar
		"AZN": 2, //  Azerbaijani manat
		"BSD": 2, //  Bahamian dollar
		"BHD": 3, //  Bahraini dinar
		"BDT": 2, //  Bangladeshi taka
		"BBD": 2, //  Barbados dollar
		"BYN": 2, //  Belarusian ruble
		"BZD": 2, //  Belize dollar
		"XOF": 0, //  West African CFA franc
		"BMD": 2, //  Bermudian dollar
		"INR": 2, //  Indian rupee
		"BTN": 2, //  Bhutanese ngultrum
		"BOB": 2, //  Bolivian boliviano
		"BOV": 2, //  Bolivian Mvdol (funds code)
		"BAM": 2, //  Bosnia and Herzegovina convertible mark
		"BWP": 2, //  Botswana pula
		"NOK": 2, //  Norwegian krone
		"BRL": 2, //  Brazilian real
		"BND": 2, //  Brunei dollar
		"BGN": 2, //  Bulgarian lev
		"BIF": 0, //  Burundian franc
		"CVE": 2, //  Cape Verdean escudo
		"KHR": 2, //  Cambodian riel
		"XAF": 0, //  Central African CFA franc
		"CAD": 2, //  Canadian dollar
		"KYD": 2, //  Cayman Islands dollar
		"CLP": 0, //  Chilean peso
		"CLF": 4, //  Unidad de Fomento (funds code)
		"CNY": 2, //  Chinese yuan
		"COP": 2, //  Colombian peso
		"COU": 2, //  Unidad de Valor Real (UVR) (funds code)[7]
		"KMF": 0, //  Comoro franc
		"CDF": 2, //  Congolese franc
		"NZD": 2, //  New Zealand dollar
		"CRC": 2, //  Costa Rican colon
		"HRK": 2, //  Croatian kuna
		"CUP": 2, //  Cuban peso
		"CUC": 2, //  Cuban convertible peso
		"ANG": 2, //  Netherlands Antillean guilder
		"CZK": 2, //  Czech koruna
		"DKK": 2, //  Danish krone
		"DJF": 0, //  Djiboutian franc
		"DOP": 2, //  Dominican peso
		"EGP": 2, //  Egyptian pound
		"SVC": 2, //  Salvadoran colón
		"ERN": 2, //  Eritrean nakfa
		"SZL": 2, //  Swazi lilangeni
		"ETB": 2, //  Ethiopian birr
		"FKP": 2, //  Falkland Islands pound
		"FJD": 2, //  Fiji dollar
		"XPF": 0, //  CFP franc
		"GMD": 2, //  Gambian dalasi
		"GEL": 2, //  Georgian lari
		"GHS": 2, //  Ghanaian cedi
		"GIP": 2, //  Gibraltar pound
		"GTQ": 2, //  Guatemalan quetzal
		"GBP": 2, //  Pound sterling
		"GNF": 0, //  Guinean franc
		"GYD": 2, //  Guyanese dollar
		"HTG": 2, //  Haitian gourde
		"HNL": 2, //  Honduran lempira
		"HKD": 2, //  Hong Kong dollar
		"HUF": 2, //  Hungarian forint
		"ISK": 0, //  Icelandic króna
		"IDR": 2, //  Indonesian rupiah
		"IRR": 2, //  Iranian rial
		"IQD": 3, //  Iraqi dinar
		"ILS": 2, //  Israeli new shekel
		"JMD": 2, //  Jamaican dollar
		"JPY": 0, //  Japanese yen
		"JOD": 3, //  Jordanian dinar
		"KZT": 2, //  Kazakhstani tenge
		"KES": 2, //  Kenyan shilling
		"KPW": 2, //  North Korean won
		"KRW": 0, //  South Korean won
		"KWD": 3, //  Kuwaiti dinar
		"KGS": 2, //  Kyrgyzstani som
		"LAK": 2, //  Lao kip
		"LBP": 2, //  Lebanese pound
		"LSL": 2, //  Lesotho loti
		"ZAR": 2, //  South African rand
		"LRD": 2, //  Liberian dollar
		"LYD": 3, //  Libyan dinar
		"CHF": 2, //  Swiss franc
		"MOP": 2, //  Macanese pataca
		"MKD": 2, //  Macedonian denar
		"MGA": 2, //  Malagasy ariary
		"MWK": 2, //  Malawian kwacha
		"MYR": 2, //  Malaysian ringgit
		"MVR": 2, //  Maldivian rufiyaa
		"MRU": 2, //  Mauritanian ouguiya
		"MUR": 2, //  Mauritian rupee
		"MXN": 2, //  Mexican peso
		"MXV": 2, //  Mexican Unidad de Inversion (UDI) (funds code)
		"MDL": 2, //  Moldovan leu
		"MNT": 2, //  Mongolian tögrög
		"MAD": 2, //  Moroccan dirham
		"MZN": 2, //  Mozambican metical
		"MMK": 2, //  Burmese kyat
		"NAD": 2, //  Namibian dollar
		"NPR": 2, //  Nepalese rupee
		"NIO": 2, //  Nicaraguan córdoba
		"NGN": 2, //  Nigerian naira
		"OMR": 3, //  Omani rial
		"PKR": 2, //  Pakistani rupee
		"PAB": 2, //  Panamanian balboa
		"PGK": 2, //  Papua New Guinean kina
		"PYG": 0, //  Paraguayan guaraní
		"PEN": 2, //  Peruvian sol
		"PHP": 2, //  Philippine peso
		"PLN": 2, //  Polish złoty
		"QAR": 2, //  Qatari riyal
		"RON": 2, //  Romanian leu
		"RUB": 2, //  Russian ruble
		"RWF": 0, //  Rwandan franc
		"SHP": 2, //  Saint Helena pound
		"WST": 2, //  Samoan tala
		"STN": 2, //  São Tomé and Príncipe dobra
		"SAR": 2, //  Saudi riyal
		"RSD": 2, //  Serbian dinar
		"SCR": 2, //  Seychelles rupee
		"SLL": 2, //  Sierra Leonean leone
		"SGD": 2, //  Singapore dollar
		"SBD": 2, //  Solomon Islands dollar
		"SOS": 2, //  Somali shilling
		"SSP": 2, //  South Sudanese pound
		"LKR": 2, //  Sri Lankan rupee
		"SDG": 2, //  Sudanese pound
		"SRD": 2, //  Surinamese dollar
		"SEK": 2, //  Swedish krona/kronor
		"CHE": 2, //  WIR Euro (complementary currency)
		"CHW": 2, //  WIR Franc (complementary currency)
		"SYP": 2, //  Syrian pound
		"TWD": 2, //  New Taiwan dollar
		"TJS": 2, //  Tajikistani somoni
		"TZS": 2, //  Tanzanian shilling
		"THB": 2, //  Thai baht
		"TOP": 2, //  Tongan paʻanga
		"TTD": 2, //  Trinidad and Tobago dollar
		"TND": 3, //  Tunisian dinar
		"TRY": 2, //  Turkish lira
		"TMT": 2, //  Turkmenistan manat
		"UGX": 0, //  Ugandan shilling
		"UAH": 2, //  Ukrainian hryvnia
		"AED": 2, //  United Arab Emirates dirham
		"USN": 2, //  United States dollar (next day) (funds code)
		"UYU": 2, //  Uruguayan peso
		"UYI": 0, //  Uruguay Peso en Unidades Indexadas (URUIURUI) (funds code)
		"UYW": 4, //  Unidad previsional[9]
		"UZS": 2, //  Uzbekistan som
		"VUV": 0, //  Vanuatu vatu
		"VES": 2, //  Venezuelan bolívar soberano
		"VND": 0, //  Vietnamese đồng
		"YER": 2, //  Yemeni rial
		"ZMW": 2, //  Zambian kwacha
		"ZWL": 2, //  Zimbabwean dollar A/10
	}
)

func FormatAsset(currencies map[string]int, cur string) string {
	asset := strings.ToUpper(string(cur))

	def, ok := currencies[asset]
	if !ok {
		return asset
	}

	if def == 0 {
		return asset
	}

	return fmt.Sprintf("%s/%d", asset, def)
}

func GetPrecision(currencies map[string]int, cur string) (int, error) {
	asset := strings.ToUpper(string(cur))

	def, ok := currencies[asset]
	if !ok {
		return 0, errors.New("missing currencies")
	}

	return def, nil
}

func GetCurrencyAndPrecisionFromAsset(currencies map[string]int, asset string) (string, int, error) {
	parts := strings.Split(asset, "/")
	if len(parts) != 2 {
		return "", 0, errors.New("invalid asset")
	}

	currency := parts[0]
	precision, err := GetPrecision(currencies, currency)
	if err != nil {
		return "", 0, err
	}

	return currency, precision, nil
}
