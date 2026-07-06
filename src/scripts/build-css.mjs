import { copyFile, mkdir } from 'node:fs/promises';
import { dirname } from 'node:path';

const files = [
  ['../node_modules/bootstrap/dist/css/bootstrap.min.css', '../public/assets/bootstrap.min.css'],
  ['../assets/css/input.css', '../public/assets/app.css'],
];

for (const [from, to] of files) {
  const src = new URL(from, import.meta.url);
  const dst = new URL(to, import.meta.url);
  await mkdir(dirname(dst.pathname), { recursive: true });
  await copyFile(src, dst);
}
