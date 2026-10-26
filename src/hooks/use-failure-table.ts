import { useState } from 'react'
import {
    exampleValues,
    withNewRow,
    withoutRow,
    withRowValue,
    type FailureTableValues,
    type RowField,
} from '@/lib/form-values'
import { parseFailureTable } from '@/lib/parsing'
import { computeIndicators, type Indicators } from '@/lib/reliability'

type Outcome = { kind: 'pending' } | { kind: 'results'; rows: Indicators[] } | { kind: 'error'; error: unknown }

const pending: Outcome = { kind: 'pending' }

function computeOutcome(values: FailureTableValues): Outcome {
    try {
        return { kind: 'results', rows: computeIndicators(parseFailureTable(values)) }
    } catch (error) {
        return { kind: 'error', error }
    }
}

function useFailureTable() {
    const [values, setValues] = useState(exampleValues)
    const [outcome, setOutcome] = useState(pending)

    const edit = (update: (values: FailureTableValues) => FailureTableValues) => {
        setValues(update)
        setOutcome(pending)
    }

    return {
        values,
        outcome,
        setInitialCount: (value: string) => edit((current) => ({ ...current, initialCount: value })),
        updateRow: (id: number, field: RowField, value: string) => edit((current) => withRowValue(current, id, field, value)),
        addRow: () => edit(withNewRow),
        removeRow: (id: number) => edit((current) => withoutRow(current, id)),
        loadExample: () => edit(exampleValues),
        compute: () => setOutcome(computeOutcome(values)),
    }
}

export { useFailureTable }
