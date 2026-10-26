type DataErrorCode =
    | 'invalidCount'
    | 'noRows'
    | 'invalidTime'
    | 'invalidFailures'
    | 'timesNotIncreasing'
    | 'failuresExceedCount'

class DataError extends Error {
    readonly code: DataErrorCode
    readonly row?: number

    constructor(code: DataErrorCode, row?: number) {
        super(row === undefined ? code : `row ${row}: ${code}`)
        this.name = 'DataError'
        this.code = code
        this.row = row
    }
}

function errorRow(error: unknown): number | undefined {
    return error instanceof DataError ? error.row : undefined
}

export { DataError, errorRow }
export type { DataErrorCode }
