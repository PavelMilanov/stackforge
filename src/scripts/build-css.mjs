import { copyFile, mkdir } from 'node:fs/promises';
import { dirname } from 'node:path';

const bootstrapSource = new URL('../node_modules/bootstrap/dist/css/bootstrap.min.css', import.meta.url);
const bootstrapTarget = new URL('../public/assets/bootstrap.min.css', import.meta.url);

await mkdir(dirname(bootstrapTarget.pathname), { recursive: true });
await copyFile(bootstrapSource, bootstrapTarget);
