import { DataError } from './errors'
import type { FailureTableValues } from './form-values'
import type { FailureTable } from './reliability'

const integerPattern = /^\d+$/
const decimalPattern = /^\d+([.,]\d+)?$/
const decimalComma = ','
const decimalPoint = '.'

function parseInteger(text: string): number | undefined {
    const trimmed = text.trim()
    return integerPattern.test(trimmed) ? Number(trimmed) : undefined
}

function parseDecimal(text: string): number | undefined {
    const trimmed = text.trim()
    return decimalPattern.test(trimmed) ? Number(trimmed.replace(decimalComma, decimalPoint)) : undefined
}

function parseRequired(
    text: string,
    parse: (text: string) => number | undefined,
    code: 'invalidTime' | 'invalidFailures',
    row: number,
): number {
    const value = parse(text)
    if (value === undefined) {
        throw new DataError(code, row)
    }
    return value
}

function parseFailureTable(values: FailureTableValues): FailureTable {
    const initialCount = parseInteger(values.initialCount)
    if (initialCount === undefined) {
        throw new DataError('invalidCount')
    }
    const table: FailureTable = { initialCount, times: [], failures: [] }
    values.rows.forEach((row, index) => {
        table.times.push(parseRequired(row.time, parseDecimal, 'invalidTime', index))
        table.failures.push(parseRequired(row.failures, parseInteger, 'invalidFailures', index))
    })
    return table
}

export { parseFailureTable }
