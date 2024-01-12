Prism.languages.numscript = {
    comment: [
        /\/\/[^\n]*/,
        /\/\*[\s\S]*?\*\//
    ],
    keyword: [
        /send|vars|remaining|to|kept|allowing|unbounded|overdraft|up|max|from/
    ],
    property: /source|destination/,
    builtin: /set_tx_meta|set_account_meta/,
    symbol: [
        /account|monetary|asset|number|string|portion/,
        /\@world/
    ],
    string: [
        /"[^"]*"/,
        /\[[A-Z0-9]+(\/\d+)? (\d+|\*)\]/
    ],
    punctuation: /\(|\{|\[|\]|\}|\)|=/,
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
        /^\d+$/
    ],
    // constant: {
    //     pattern: /@[A-Za-z0-9:]+/,
    //     greedy: true
    // }
}