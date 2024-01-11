Prism.languages.numscript = {
    comment: [
        /\/\/[^\n]*/,
        /\/\*[\s\S]*?\*\//
    ],
    keyword: [
        /send|vars|source|destination|remaining|to|kept|allowing|unbounded|overdraft|up|max|from/
    ],
    builtin: /set_tx_meta|set_account_meta/,
    symbol: /account|monetary|asset|number|string|portion/,
    string: [
        /"[^"]*"/,
        /\[[A-Z0-9]+(\/\d+)? \d+\]/,
        {
            pattern: /@[A-Za-z0-9:]+/,
            greedy: true
        }
    ],
    punctuation: /\(\{\[\]\}\)=/,
    variable: /\$[a-zA-Z][A-Za-z0-9]*/,
    number: [
        {
            pattern: /\d+\/\d+/,
            greedy: true
        }, 
        {
            pattern: /\d+%/,
            greedy: true
        },
        /\d+/
    ]
}