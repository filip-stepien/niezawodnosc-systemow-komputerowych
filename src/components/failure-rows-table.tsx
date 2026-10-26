import { XIcon } from 'lucide-react'
import { ScrollableTable } from '@/components/scrollable-table'
import { Tex } from '@/components/tex'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { useScrollToEndOnAdd } from '@/hooks/use-scroll-to-end-on-add'
import type { RowField, RowValues } from '@/lib/form-values'
import { cn } from '@/lib/utils'

type FailureRowProps = {
    row: RowValues
    number: number
    invalid: boolean
    onChange: (id: number, field: RowField, value: string) => void
    onRemove: (id: number) => void
}

type FailureRowsTableProps = {
    rows: RowValues[]
    errorRow?: number
    onChange: (id: number, field: RowField, value: string) => void
    onRemove: (id: number) => void
}

function FailureRow({ row, number, invalid, onChange, onRemove }: FailureRowProps) {
    const inputClass = cn('text-right tabular-nums', invalid && 'border-destructive')
    return (
        <TableRow>
            <TableCell className="text-muted-foreground tabular-nums">{number}</TableCell>
            <TableCell>
                <Input
                    className={inputClass}
                    inputMode="decimal"
                    value={row.time}
                    onChange={(event) => onChange(row.id, 'time', event.target.value)}
                />
            </TableCell>
            <TableCell>
                <Input
                    className={inputClass}
                    inputMode="numeric"
                    value={row.failures}
                    onChange={(event) => onChange(row.id, 'failures', event.target.value)}
                />
            </TableCell>
            <TableCell className="text-right">
                <Button variant="ghost" size="icon-sm" onClick={() => onRemove(row.id)}>
                    <XIcon />
                </Button>
            </TableCell>
        </TableRow>
    )
}

function FailureRowsTable({ rows, errorRow, onChange, onRemove }: FailureRowsTableProps) {
    const tableRef = useScrollToEndOnAdd(rows.length)
    return (
        <ScrollableTable ref={tableRef}>
            <TableHeader>
                <TableRow>
                    <TableHead className="w-10">Lp.</TableHead>
                    <TableHead>
                        Czas <Tex>t</Tex>
                    </TableHead>
                    <TableHead>
                        Uszkodzenia <Tex>m(t)</Tex>
                    </TableHead>
                    <TableHead className="w-10" />
                </TableRow>
            </TableHeader>
            <TableBody>
                {rows.map((row, index) => (
                    <FailureRow
                        key={row.id}
                        row={row}
                        number={index + 1}
                        invalid={index === errorRow}
                        onChange={onChange}
                        onRemove={onRemove}
                    />
                ))}
            </TableBody>
        </ScrollableTable>
    )
}

export { FailureRowsTable }
