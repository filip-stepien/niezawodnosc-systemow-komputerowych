import type { ReactNode } from 'react'
import { Tex } from '@/components/tex'
import { DataError, type DataErrorCode } from '@/lib/errors'

type ErrorDescriptionProps = {
    error: unknown
}

const messages: Record<DataErrorCode, ReactNode> = {
    invalidCount: (
        <>
            Liczność próby <Tex>n(0)</Tex> musi być dodatnią liczbą całkowitą.
        </>
    ),
    noRows: 'Tabela uszkodzeń jest pusta.',
    invalidTime: (
        <>
            czas <Tex>t</Tex> musi być liczbą dodatnią.
        </>
    ),
    invalidFailures: (
        <>
            liczba uszkodzeń <Tex>m(t)</Tex> musi być nieujemną liczbą całkowitą.
        </>
    ),
    timesNotIncreasing: (
        <>
            czasy <Tex>t</Tex> muszą być dodatnie i rosnąć z wiersza na wiersz.
        </>
    ),
    failuresExceedCount: (
        <>
            Łączna liczba uszkodzeń przekracza liczność próby <Tex>n(0)</Tex>.
        </>
    ),
}

function ErrorDescription({ error }: ErrorDescriptionProps) {
    if (!(error instanceof DataError)) {
        return <>Nieoczekiwany błąd: {String(error)}</>
    }
    if (error.row === undefined) {
        return messages[error.code]
    }
    return (
        <>
            Wiersz {error.row + 1}: {messages[error.code]}
        </>
    )
}

export { ErrorDescription }
