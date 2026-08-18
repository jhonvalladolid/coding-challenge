import { readFileSync } from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';
import swaggerUi from 'swagger-ui-express';
import { load } from 'js-yaml';

const specPath = path.join(path.dirname(fileURLToPath(import.meta.url)), '../../docs/openapi.yaml');
const spec = load(readFileSync(specPath, 'utf8'));

export function registerDocs(app) {
  app.use('/docs', swaggerUi.serve, swaggerUi.setup(spec, {
    customSiteTitle: 'Node API — Swagger',
  }));
}
