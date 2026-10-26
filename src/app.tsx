import { FailureTableForm } from '@/components/failure-table-form'
import { ResultsTable } from '@/components/results-table'
import { Tex } from '@/components/tex'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { useFailureTable } from '@/hooks/use-failure-table'
import { errorRow } from '@/lib/errors'

function App() {
    const table = useFailureTable()
    const failure = table.outcome.kind === 'error' ? table.outcome.error : undefined

    return (
        <div className="min-h-screen bg-muted/40">
            <div className="mx-auto flex max-w-6xl flex-col gap-6 px-4 py-8">
                <header className="flex flex-col gap-1">
                    <h1 className="font-heading text-2xl font-semibold">Empiryczne wskaźniki niezawodności</h1>
                    <p className="text-muted-foreground">
                        <Tex>F^*(t)</Tex>, <Tex>R^*(t)</Tex>, <Tex>f^*(t)</Tex> i <Tex>\lambda^*(t)</Tex> wyznaczane z
                        tabeli uszkodzeń
                    </p>
                </header>
                <main className="grid items-start gap-6 lg:grid-cols-[minmax(0,28rem)_minmax(0,1fr)]">
                    <Card>
                        <CardHeader>
                            <CardTitle>Dane</CardTitle>
                        </CardHeader>
                        <CardContent>
                            <FailureTableForm
                                values={table.values}
                                error={failure}
                                errorRow={errorRow(failure)}
                                onInitialCountChange={table.setInitialCount}
                                onRowChange={table.updateRow}
                                onRowAdd={table.addRow}
                                onRowRemove={table.removeRow}
                                onLoadExample={table.loadExample}
                                onCompute={table.compute}
                            />
                        </CardContent>
                    </Card>
                    <Card>
                        <CardHeader>
                            <CardTitle>Wyniki</CardTitle>
                        </CardHeader>
                        <CardContent>
                            <ResultsTable rows={table.outcome.kind === 'results' ? table.outcome.rows : []} />
                        </CardContent>
                    </Card>
                </main>
            </div>
        </div>
    )
}

export default App
