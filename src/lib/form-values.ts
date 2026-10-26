import { bookExample } from '@/constants/example'

type RowValues = {
    id: number
    time: string
    failures: string
}

type FailureTableValues = {
    initialCount: string
    rows: RowValues[]
}

type RowField = 'time' | 'failures'

let nextRowId = 0

function createRow(time = '', failures = ''): RowValues {
    return { id: nextRowId++, time, failures }
}

function exampleValues(): FailureTableValues {
    return {
        initialCount: String(bookExample.initialCount),
        rows: bookExample.rows.map(([time, failures]) => createRow(String(time), String(failures))),
    }
}

function withRowValue(values: FailureTableValues, id: number, field: RowField, value: string): FailureTableValues {
    return { ...values, rows: values.rows.map((row) => (row.id === id ? { ...row, [field]: value } : row)) }
}

function withoutRow(values: FailureTableValues, id: number): FailureTableValues {
    return { ...values, rows: values.rows.filter((row) => row.id !== id) }
}

function withNewRow(values: FailureTableValues): FailureTableValues {
    return { ...values, rows: [...values.rows, createRow()] }
}

export { exampleValues, withNewRow, withoutRow, withRowValue }
export type { FailureTableValues, RowField, RowValues }
