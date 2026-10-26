import { DataError } from './errors'
import type { FailureTable } from './reliability'

function sum(values: number[]): number {
    return values.reduce((total, value) => total + value, 0)
}

function validateRow(table: FailureTable, time: number, row: number): void {
    const previousTime = row === 0 ? 0 : table.times[row - 1]
    if (!Number.isFinite(time) || time <= previousTime) {
        throw new DataError('timesNotIncreasing', row)
    }
    const failures = table.failures[row]
    if (!Number.isInteger(failures) || failures < 0) {
        throw new DataError('invalidFailures', row)
    }
}

function validateFailureTable(table: FailureTable): void {
    if (!Number.isInteger(table.initialCount) || table.initialCount <= 0) {
        throw new DataError('invalidCount')
    }
    if (table.times.length === 0) {
        throw new DataError('noRows')
    }
    table.times.forEach((time, row) => validateRow(table, time, row))
    if (sum(table.failures) > table.initialCount) {
        throw new DataError('failuresExceedCount')
    }
}

export { validateFailureTable }
