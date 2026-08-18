import express from 'express';
import { env } from './config/env.js';
import { errorHandler } from './middlewares/error.middleware.js';
import { notFound } from './middlewares/not-found.middleware.js';
import { requestId } from './middlewares/request-id.middleware.js';
import statisticsRoutes from './modules/statistics/statistics.routes.js';
import { registerDocs } from './config/swagger.js';
import { success } from './utils/response.util.js';

const app = express();

app.disable('x-powered-by');
app.use(requestId);
app.use(express.json({ limit: env.jsonBodyLimit }));

app.get('/health', (req, res) => {
  res.status(200).json(
    success({
      status: 'ok',
      service: 'node-api',
    }),
  );
});

registerDocs(app);
app.use('/api/v1/statistics', statisticsRoutes);
app.use(notFound);
app.use(errorHandler);

export default app;
