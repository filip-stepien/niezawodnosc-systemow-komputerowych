import type { ReactNode, Ref } from 'react'
import { Table } from '@/components/ui/table'
import { cn } from '@/lib/utils'

type ScrollableTableProps = {
    ref?: Ref<HTMLDivElement>
    children: ReactNode
}

const containerClass = cn(
    '*:data-[slot=table-container]:max-h-[26rem]',
    '*:data-[slot=table-container]:overflow-y-auto',
    '*:data-[slot=table-container]:scrollbar-subtle',
)

const headerClass = cn(
    '[&_thead]:sticky',
    '[&_thead]:top-0',
    '[&_thead]:bg-card',
    '[&_thead_tr]:border-b-0',
    '[&_th]:shadow-[inset_0_-1px_0_var(--border)]',
)

const rowClass = '[&_tr]:hover:bg-transparent'

function ScrollableTable({ ref, children }: ScrollableTableProps) {
    return (
        <div ref={ref} className={cn(containerClass, headerClass, rowClass)}>
            <Table>{children}</Table>
        </div>
    )
}

export { ScrollableTable }
