import katex from 'katex'
import 'katex/dist/katex.min.css'
import { useMemo } from 'react'

type TexProps = {
    children: string
}

function Tex({ children }: TexProps) {
    const html = useMemo(() => katex.renderToString(children, { output: 'html', throwOnError: false }), [children])
    return <span dangerouslySetInnerHTML={{ __html: html }} />
}

export { Tex }
