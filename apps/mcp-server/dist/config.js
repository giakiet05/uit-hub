import dotenv from 'dotenv';
import { dirname, resolve } from 'path';
import { fileURLToPath } from 'url';
const __dirname = dirname(fileURLToPath(import.meta.url));
dotenv.config({ path: resolve(__dirname, '../../.env') }); // Dùng chung .env của monorepo (nếu có) hoặc .env trong apps/mcp-server
export const config = {
    // Fake UIT Server Base URL
    apiBaseUrl: process.env.API_BASE_URL || 'http://localhost:8080/api/v1',
};
