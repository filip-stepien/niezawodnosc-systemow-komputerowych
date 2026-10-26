import { useEffect, useRef } from 'react'

const scrollContainerSelector = '[data-slot=table-container]'

function useScrollToEndOnAdd(count: number) {
    const ref = useRef<HTMLDivElement>(null)
    const previousCount = useRef(count)

    useEffect(() => {
        const container = ref.current?.querySelector(scrollContainerSelector)
        if (container && count === previousCount.current + 1) {
            container.scrollTo({ top: container.scrollHeight, behavior: 'smooth' })
        }
        previousCount.current = count
    }, [count])

    return ref
}

export { useScrollToEndOnAdd }
