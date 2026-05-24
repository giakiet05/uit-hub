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

  server.tool(
    "room_plan_meeting",
    "Đề xuất phòng họp phù hợp từ danh sách phòng trống.",
    {
      token: z.string().describe("Bearer token"),
      ...QueryRoomsAvailabilitySchema,
      capacity: z.number().optional().describe("Sức chứa cần thiết"),
      building: z.string().optional().describe("Tòa nhà (VD: A, B, E)"),
      equipment: z.array(z.string()).optional().describe("Danh sách trang thiết bị cần thiết (VD: projector, microphone)")
    },
    async ({ token, date, start, end, capacity, building, equipment }) => {
      const result = await callFakeServer(
        "room_get_availability", 
        "GET", 
        "/rooms/availability", 
        token, 
        undefined, 
        { date, start, end }
      );

      if (!result.ok) {
        return formatMcpResponse(result);
      }

      let rooms: any[] = Array.isArray(result.data) ? result.data : (result.data?.rooms || []);

      if (capacity) {
        rooms = rooms.filter((r: any) => r.capacity && r.capacity >= capacity);
      }
      if (building) {
        rooms = rooms.filter((r: any) => r.building === building);
      }
      if (equipment && equipment.length > 0) {
        rooms = rooms.filter((r: any) => {
          const roomEqs = r.equipment || [];
          return equipment.every((e: string) => roomEqs.includes(e));
        });
      }

      return formatMcpResponse({
        ok: true,
        tool: "room_plan_meeting",
        data: {
          requested_date: date,
          start,
          end,
          suggested_rooms: rooms
        }
      });
    }
  );
}
