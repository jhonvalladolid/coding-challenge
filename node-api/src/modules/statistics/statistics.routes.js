import { Router } from 'express';
import { methodNotAllowed } from '../../middlewares/not-found.middleware.js';
import { createStatistics } from './statistics.controller.js';
import {
  assertJsonContentType,
  validateStatisticsBody,
} from './statistics.validator.js';

const router = Router();

router.post('/', assertJsonContentType, validateStatisticsBody, createStatistics);
router.all('/', methodNotAllowed);

export default router;
