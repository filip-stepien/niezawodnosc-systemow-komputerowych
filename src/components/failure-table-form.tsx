import { CircleAlertIcon, PlusIcon } from 'lucide-react'
import { ErrorDescription } from '@/components/error-description'
import { FailureRowsTable } from '@/components/failure-rows-table'
import { Tex } from '@/components/tex'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { Field, FieldDescription, FieldGroup, FieldLabel, FieldTitle } from '@/components/ui/field'
import { Input } from '@/components/ui/input'
import type { FailureTableValues, RowField } from '@/lib/form-values'

type FailureTableFormProps = {
    values: FailureTableValues
    error?: unknown
    errorRow?: number
    onInitialCountChange: (value: string) => void
    onRowChange: (id: number, field: RowField, value: string) => void
    onRowAdd: () => void
    onRowRemove: (id: number) => void
    onLoadExample: () => void
    onCompute: () => void
}

function FailureTableForm(props: FailureTableFormProps) {
    return (
        <form
            noValidate
            className="flex flex-col gap-6"
            onSubmit={(event) => {
                event.preventDefault()
                props.onCompute()
            }}
        >
            <FieldGroup>
                <Field orientation="horizontal">
                    <FieldLabel htmlFor="initial-count">
                        <span>
                            Liczność próby <Tex>n(0)</Tex>
                        </span>
                    </FieldLabel>
                    <Input
                        id="initial-count"
                        className="w-32 text-right tabular-nums"
                        inputMode="numeric"
                        value={props.values.initialCount}
                        onChange={(event) => props.onInitialCountChange(event.target.value)}
                    />
                </Field>
                <Field>
                    <FieldTitle>Tabela uszkodzeń</FieldTitle>
                    <FieldDescription>
                        <Tex>m(t)</Tex> to liczba elementów, które uszkodziły się od poprzedniej chwili do chwili <Tex>t</Tex>.
                    </FieldDescription>
                    <FailureRowsTable
                        rows={props.values.rows}
                        errorRow={props.errorRow}
                        onChange={props.onRowChange}
                        onRemove={props.onRowRemove}
                    />
                    <Button variant="outline" onClick={props.onRowAdd}>
                        <PlusIcon data-icon="inline-start" />
                        Dodaj wiersz
                    </Button>
                </Field>
            </FieldGroup>
            {props.error !== undefined && (
                <Alert variant="destructive">
                    <CircleAlertIcon />
                    <AlertDescription>
                        <ErrorDescription error={props.error} />
                    </AlertDescription>
                </Alert>
            )}
            <div className="flex gap-2">
                <Button variant="outline" onClick={props.onLoadExample}>
                    Wczytaj przykład
                </Button>
                <Button type="submit" className="flex-1">
                    Oblicz
                </Button>
            </div>
        </form>
    )
}

export { FailureTableForm }
