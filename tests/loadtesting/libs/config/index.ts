export const loadConfig = (): Index => {
    return {
        ledgerUrl: __ENV.LEDGER_URL,
    };
};

export interface Index {
    ledgerUrl: string;
}
