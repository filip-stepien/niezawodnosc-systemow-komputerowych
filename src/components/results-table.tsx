import { ScrollableTable } from '@/components/scrollable-table'
import { Tex } from '@/components/tex'
import { Empty, EmptyDescription, EmptyHeader, EmptyTitle } from '@/components/ui/empty'
import { TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { formatIndicator, formatTime } from '@/lib/format'
import type { Indicators } from '@/lib/reliability'

type ResultColumn = {
    header: string
    value: (row: Indicators) => string
}

const columns: ResultColumn[] = [
    { header: 't', value: (row) => formatTime(row.time) },
    { header: 'F^*(t)', value: (row) => formatIndicator(row.unreliability) },
    { header: 'R^*(t)', value: (row) => formatIndicator(row.reliability) },
    { header: 'f^*(t)', value: (row) => formatIndicator(row.density) },
    { header: '\\lambda^*(t)', value: (row) => formatIndicator(row.hazardRate) },
]

function ResultsTable({ rows }: { rows: Indicators[] }) {
    if (rows.length === 0) {
        return (
            <Empty>
                <EmptyHeader>
                    <EmptyTitle>Brak wyników</EmptyTitle>
                    <EmptyDescription>Wyniki pojawią się po kliknięciu "Oblicz".</EmptyDescription>
                </EmptyHeader>
            </Empty>
        )
    }
    return (
        <ScrollableTable>
            <TableHeader>
                <TableRow>
                    {columns.map((column) => (
                        <TableHead key={column.header} className="text-right">
                            <Tex>{column.header}</Tex>
                        </TableHead>
                    ))}
                </TableRow>
            </TableHeader>
            <TableBody>
                {rows.map((row) => (
                    <TableRow key={row.time}>
                        {columns.map((column) => (
                            <TableCell key={column.header} className="text-right tabular-nums">
                                {column.value(row)}
                            </TableCell>
                        ))}
                    </TableRow>
                ))}
            </TableBody>
        </ScrollableTable>
    )
}

export { ResultsTable }
