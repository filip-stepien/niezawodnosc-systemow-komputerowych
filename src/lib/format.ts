const locale = 'pl-PL'
const indicatorDigits = 3
const maxTimeDigits = 6
const undefinedValue = '-'

const timeFormat = new Intl.NumberFormat(locale, {
    maximumFractionDigits: maxTimeDigits,
    useGrouping: false,
})

const indicatorFormat = new Intl.NumberFormat(locale, {
    minimumFractionDigits: indicatorDigits,
    maximumFractionDigits: indicatorDigits,
    useGrouping: false,
})

function formatTime(value: number): string {
    return timeFormat.format(value)
}

function formatIndicator(value: number): string {
    return Number.isNaN(value) ? undefinedValue : indicatorFormat.format(value)
}

export { formatIndicator, formatTime }
