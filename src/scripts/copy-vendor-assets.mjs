import { copyFile, mkdir } from 'node:fs/promises';
import { dirname } from 'node:path';

const files = [
  ['../node_modules/htmx.org/dist/htmx.min.js', '../public/assets/htmx.min.js'],
  ['../node_modules/alpinejs/dist/cdn.min.js', '../public/assets/alpine.min.js'],
];

for (const [from, to] of files) {
  const src = new URL(from, import.meta.url);
  const dst = new URL(to, import.meta.url);
  await mkdir(dirname(dst.pathname), { recursive: true });
  await copyFile(src, dst);
}
