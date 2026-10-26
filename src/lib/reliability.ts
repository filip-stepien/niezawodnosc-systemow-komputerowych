import { validateFailureTable } from './validation'

type FailureTable = {
    initialCount: number
    times: number[]
    failures: number[]
}

type Indicators = {
    time: number
    unreliability: number
    reliability: number
    density: number
    hazardRate: number
}

type Interval = {
    time: number
    length: number
    failures: number
    survivorsBefore: number
}

function toIntervals(table: FailureTable): Interval[] {
    const intervals: Interval[] = []
    let survivorsBefore = table.initialCount
    let previousTime = 0
    table.times.forEach((time, row) => {
        const failures = table.failures[row]
        intervals.push({ time, length: time - previousTime, failures, survivorsBefore })
        survivorsBefore -= failures
        previousTime = time
    })
    return intervals
}

function hazardRate(interval: Interval): number {
    if (interval.survivorsBefore === 0) {
        return Number.NaN
    }
    return interval.failures / (interval.survivorsBefore * interval.length)
}

function indicatorsFor(initialCount: number, interval: Interval): Indicators {
    const survivorsAfter = interval.survivorsBefore - interval.failures
    return {
        time: interval.time,
        unreliability: (initialCount - survivorsAfter) / initialCount,
        reliability: survivorsAfter / initialCount,
        density: interval.failures / (initialCount * interval.length),
        hazardRate: hazardRate(interval),
    }
}

function computeIndicators(table: FailureTable): Indicators[] {
    validateFailureTable(table)
    return toIntervals(table).map((interval) => indicatorsFor(table.initialCount, interval))
}

export { computeIndicators }
export type { FailureTable, Indicators }
