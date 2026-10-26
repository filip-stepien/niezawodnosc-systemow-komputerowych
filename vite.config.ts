import tailwindcss from '@tailwindcss/vite'
import react from '@vitejs/plugin-react'
import { fileURLToPath } from 'node:url'
import { defineConfig, type Plugin } from 'vite'
import { viteSingleFile } from 'vite-plugin-singlefile'

const katexStylesheet = /katex(\.min)?\.css$/
const legacyFontSource = /,\s*url\([^)]+\.(woff|ttf)\)\s*format\("(woff|truetype)"\)/g

function katexWoff2Only(): Plugin {
    return {
        name: 'katex-woff2-only',
        enforce: 'pre',
        transform(code, id) {
            return katexStylesheet.test(id) ? code.replace(legacyFontSource, '') : undefined
        },
    }
}

// https://vite.dev/config/
export default defineConfig({
    plugins: [katexWoff2Only(), react(), tailwindcss(), viteSingleFile()],
    resolve: {
        alias: {
            '@': fileURLToPath(new URL('./src', import.meta.url)),
        },
    },
})
