package stripe

import "github.com/formancehq/payments/cmd/connectors/internal/connectors/currency"

var (
	// c.f. https://stripe.com/docs/currencies#zero-decimal
	supportedCurrenciesWithDecimal = map[string]int{
		"USD": currency.ISO4217Currencies["USD"], //  United States dollar
		"AED": currency.ISO4217Currencies["AED"], //  United Arab Emirates dirham
		"AFN": currency.ISO4217Currencies["AFN"], //  Afghan afghani
		"ALL": currency.ISO4217Currencies["ALL"], //  Albanian lek
		"AMD": currency.ISO4217Currencies["AMD"], //  Armenian dram
		"ANG": currency.ISO4217Currencies["ANG"], //  Netherlands Antillean guilder
		"AOA": currency.ISO4217Currencies["AOA"], //  Angolan kwanza
		"ARS": currency.ISO4217Currencies["ARS"], //  Argentine peso
		"AUD": currency.ISO4217Currencies["AUD"], //  Australian dollar
		"AWG": currency.ISO4217Currencies["AWG"], //  Aruban florin
		"AZN": currency.ISO4217Currencies["AZN"], //  Azerbaijani manat
		"BAM": currency.ISO4217Currencies["BAM"], //  Bosnia and Herzegovina convertible mark
		"BBD": currency.ISO4217Currencies["BBD"], //  Barbados dollar
		"BDT": currency.ISO4217Currencies["BDT"], //  Bangladeshi taka
		"BGN": currency.ISO4217Currencies["BGN"], //  Bulgarian lev
		"BIF": currency.ISO4217Currencies["BIF"], //  Burundian franc
		"BMD": currency.ISO4217Currencies["BMD"], //  Bermudian dollar
		"BND": currency.ISO4217Currencies["BND"], //  Brunei dollar
		"BOB": currency.ISO4217Currencies["BOB"], //  Bolivian boliviano
		"BRL": currency.ISO4217Currencies["BRL"], //  Brazilian real
		"BSD": currency.ISO4217Currencies["BSD"], //  Bahamian dollar
		"BWP": currency.ISO4217Currencies["BWP"], //  Botswana pula
		"BYN": currency.ISO4217Currencies["BYN"], //  Belarusian ruble
		"BZD": currency.ISO4217Currencies["BZD"], //  Belize dollar
		"CAD": currency.ISO4217Currencies["CAD"], //  Canadian dollar
		"CDF": currency.ISO4217Currencies["CDF"], //  Congolese franc
		"CHF": currency.ISO4217Currencies["CHF"], //  Swiss franc
		"CLP": currency.ISO4217Currencies["CLP"], //  Chilean peso
		"CNY": currency.ISO4217Currencies["CNY"], //  Chinese yuan
		"COP": currency.ISO4217Currencies["COP"], //  Colombian peso
		"CRC": currency.ISO4217Currencies["CRC"], //  Costa Rican colon
		"CVE": currency.ISO4217Currencies["CVE"], //  Cape Verdean escudo
		"CZK": currency.ISO4217Currencies["CZK"], //  Czech koruna
		"DJF": currency.ISO4217Currencies["DJF"], //  Djiboutian franc
		"DKK": currency.ISO4217Currencies["DKK"], //  Danish krone
		"DOP": currency.ISO4217Currencies["DOP"], //  Dominican peso
		"DZD": currency.ISO4217Currencies["DZD"], //  Algerian dinar
		"EGP": currency.ISO4217Currencies["EGP"], //  Egyptian pound
		"ETB": currency.ISO4217Currencies["ETB"], //  Ethiopian birr
		"EUR": currency.ISO4217Currencies["EUR"], //  Euro
		"FJD": currency.ISO4217Currencies["FJD"], //  Fiji dollar
		"FKP": currency.ISO4217Currencies["FKP"], //  Falkland Islands pound
		"GBP": currency.ISO4217Currencies["GBP"], //  Pound sterling
		"GEL": currency.ISO4217Currencies["GEL"], //  Georgian lari
		"GIP": currency.ISO4217Currencies["GIP"], //  Gibraltar pound
		"GMD": currency.ISO4217Currencies["GMD"], //  Gambian dalasi
		"GNF": currency.ISO4217Currencies["GNF"], //  Guinean franc
		"GTQ": currency.ISO4217Currencies["GTQ"], //  Guatemalan quetzal
		"GYD": currency.ISO4217Currencies["GYD"], //  Guyanese dollar
		"HKD": currency.ISO4217Currencies["HKD"], //  Hong Kong dollar
		"HNL": currency.ISO4217Currencies["HNL"], //  Honduran lempira
		"HTG": currency.ISO4217Currencies["HTG"], //  Haitian gourde
		"IDR": currency.ISO4217Currencies["IDR"], //  Indonesian rupiah
		"ILS": currency.ISO4217Currencies["ILS"], //  Israeli new shekel
		"INR": currency.ISO4217Currencies["INR"], //  Indian rupee
		"JMD": currency.ISO4217Currencies["JMD"], //  Jamaican dollar
		"JPY": currency.ISO4217Currencies["JPY"], //  Japanese yen
		"KES": currency.ISO4217Currencies["KES"], //  Kenyan shilling
		"KGS": currency.ISO4217Currencies["KGS"], //  Kyrgyzstani som
		"KHR": currency.ISO4217Currencies["KHR"], //  Cambodian riel
		"KMF": currency.ISO4217Currencies["KMF"], //  Comoro franc
		"KRW": currency.ISO4217Currencies["KRW"], //  South Korean won
		"KYD": currency.ISO4217Currencies["KYD"], //  Cayman Islands dollar
		"KZT": currency.ISO4217Currencies["KZT"], //  Kazakhstani tenge
		"LAK": currency.ISO4217Currencies["LAK"], //  Lao kip
		"LBP": currency.ISO4217Currencies["LBP"], //  Lebanese pound
		"LKR": currency.ISO4217Currencies["LKR"], //  Sri Lankan rupee
		"LRD": currency.ISO4217Currencies["LRD"], //  Liberian dollar
		"LSL": currency.ISO4217Currencies["LSL"], //  Lesotho loti
		"MAD": currency.ISO4217Currencies["MAD"], //  Moroccan dirham
		"MDL": currency.ISO4217Currencies["MDL"], //  Moldovan leu
		"MKD": currency.ISO4217Currencies["MKD"], //  Macedonian denar
		"MMK": currency.ISO4217Currencies["MMK"], //  Burmese kyat
		"MNT": currency.ISO4217Currencies["MNT"], //  Mongolian tögrög
		"MOP": currency.ISO4217Currencies["MOP"], //  Macanese pataca
		"MUR": currency.ISO4217Currencies["MUR"], //  Mauritian rupee
		"MVR": currency.ISO4217Currencies["MVR"], //  Maldivian rufiyaa
		"MWK": currency.ISO4217Currencies["MWK"], //  Malawian kwacha
		"MXN": currency.ISO4217Currencies["MXN"], //  Mexican peso
		"MYR": currency.ISO4217Currencies["MYR"], //  Malaysian ringgit
		"MZN": currency.ISO4217Currencies["MZN"], //  Mozambican metical
		"NAD": currency.ISO4217Currencies["NAD"], //  Namibian dollar
		"NGN": currency.ISO4217Currencies["NGN"], //  Nigerian naira
		"NIO": currency.ISO4217Currencies["NIO"], //  Nicaraguan córdoba
		"NOK": currency.ISO4217Currencies["NOK"], //  Norwegian krone
		"NPR": currency.ISO4217Currencies["NPR"], //  Nepalese rupee
		"NZD": currency.ISO4217Currencies["NZD"], //  New Zealand dollar
		"PAB": currency.ISO4217Currencies["PAB"], //  Panamanian balboa
		"PEN": currency.ISO4217Currencies["PEN"], //  Peruvian sol
		"PGK": currency.ISO4217Currencies["PGK"], //  Papua New Guinean kina
		"PHP": currency.ISO4217Currencies["PHP"], //  Philippine peso
		"PKR": currency.ISO4217Currencies["PKR"], //  Pakistani rupee
		"PLN": currency.ISO4217Currencies["PLN"], //  Polish złoty
		"PYG": currency.ISO4217Currencies["PYG"], //  Paraguayan guaraní
		"QAR": currency.ISO4217Currencies["QAR"], //  Qatari riyal
		"RON": currency.ISO4217Currencies["RON"], //  Romanian leu
		"RSD": currency.ISO4217Currencies["RSD"], //  Serbian dinar
		"RUB": currency.ISO4217Currencies["RUB"], //  Russian ruble
		"RWF": currency.ISO4217Currencies["RWF"], //  Rwandan franc
		"SAR": currency.ISO4217Currencies["SAR"], //  Saudi riyal
		"SBD": currency.ISO4217Currencies["SBD"], //  Solomon Islands dollar
		"SCR": currency.ISO4217Currencies["SCR"], //  Seychelles rupee
		"SEK": currency.ISO4217Currencies["SEK"], //  Swedish krona/kronor
		"SGD": currency.ISO4217Currencies["SGD"], //  Singapore dollar
		"SHP": currency.ISO4217Currencies["SHP"], //  Saint Helena pound
		"SOS": currency.ISO4217Currencies["SOS"], //  Somali shilling
		"SRD": currency.ISO4217Currencies["SRD"], //  Surinamese dollar
		"SZL": currency.ISO4217Currencies["SZL"], //  Swazi lilangeni
		"THB": currency.ISO4217Currencies["THB"], //  Thai baht
		"TJS": currency.ISO4217Currencies["TJS"], //  Tajikistani somoni
		"TOP": currency.ISO4217Currencies["TOP"], //  Tongan paʻanga
		"TRY": currency.ISO4217Currencies["TRY"], //  Turkish lira
		"TTD": currency.ISO4217Currencies["TTD"], //  Trinidad and Tobago dollar
		"TZS": currency.ISO4217Currencies["TZS"], //  Tanzanian shilling
		"UAH": currency.ISO4217Currencies["UAH"], //  Ukrainian hryvnia
		"UYU": currency.ISO4217Currencies["UYU"], //  Uruguayan peso
		"UZS": currency.ISO4217Currencies["UZS"], //  Uzbekistan som
		"VND": currency.ISO4217Currencies["VND"], //  Vietnamese đồng
		"VUV": currency.ISO4217Currencies["VUV"], //  Vanuatu vatu
		"WST": currency.ISO4217Currencies["WST"], //  Samoan tala
		"XAF": currency.ISO4217Currencies["XAF"], //  Central African CFA franc
		"XCD": currency.ISO4217Currencies["XCD"], //  East Caribbean dollar
		"XOF": currency.ISO4217Currencies["XOF"], //  West African CFA franc
		"XPF": currency.ISO4217Currencies["XPF"], //  CFP franc
		"YER": currency.ISO4217Currencies["YER"], //  Yemeni rial
		"ZAR": currency.ISO4217Currencies["ZAR"], //  South African rand
		"ZMW": currency.ISO4217Currencies["ZMW"], //  Zambian kwacha

		// Unsupported currencies
		// The following currencies are not in the ISO 4217 standard,
		//so let's not handle them for now.
		// "SLE": 2 //  Sierra Leonean leone
		// "STD": 2 //  São Tomé and Príncipe dobra

		// The following currencies have not the same decimals in Stripe compared
		// to ISO 4217 standard, so let's not handle them for now.
		// "MGA": 2, //  Malagasy ariary

		// The following currencies are 3 decimals currencies, but in order
		// to use them with Stripe, it requires the last digit to be 0.
		// Let's not handle them for now.
		// "BHD": 3, //  Bahraini dinar
		// "JOD": 3, //  Jordanian dinar
		// "KWD": 3, //  Kuwaiti dinar
		// "OMR": 3, //  Omani rial
		// "TND": 3, //  Tunisian dinar

		// The following currencies are apecial cases in stripe API (cf link above)
		// let's not handle them for now.
		// "ISK": 0, //  Icelandic króna
		// "HUF": 2, //  Hungarian forint
		// "UGX": 0, //  Ugandan shilling
		// "TWD": 2 // 	New Taiwan dollar
	}
)
