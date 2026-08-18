import { success } from '../../utils/response.util.js';
import { computeStatistics } from './statistics.service.js';

export function createStatistics(req, res) {
  const { q, r } = req.body.matrices;
  const data = computeStatistics(q, r);
  res.status(200).json(success(data));
}
