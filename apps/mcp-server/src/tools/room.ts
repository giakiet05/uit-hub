import { McpServer } from '@modelcontextprotocol/sdk/server/mcp.js';
import { z } from 'zod';
import { callFakeServer, formatMcpResponse } from '../utils/api.js';
import { QueryRoomsAvailabilitySchema } from '../utils/types.js';

export function registerRoomTools(server: McpServer) {
  server.tool(
    "room_get_availability",
    "Tra cứu phòng trống.",
    {
      token: z.string().describe("Bearer token"),
      ...QueryRoomsAvailabilitySchema
    },
    async ({ token, date, start, end }) => {
      const result = await callFakeServer(
        "room_get_availability", 
        "GET", 
        "/rooms/availability", 
        token, 
        undefined, 
        { date, start, end }
      );
      return formatMcpResponse(result);
    }
  );
}
