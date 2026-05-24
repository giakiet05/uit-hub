import { registerAuthTools } from './src/tools/auth.js';
import { registerStudentTools } from './src/tools/student.js';
import { registerRoomTools } from './src/tools/room.js';
import { config } from './src/config.js';

async function test() {
  console.log('Testing MCP to Fake-UIT-Server on', config.apiBaseUrl);

  const tools: Record<string, Function> = {};

  const mockServer = {
    tool: (name: string, desc: string, schema: any, handler: Function) => {
      tools[name] = handler;
    }
  } as any;

  registerAuthTools(mockServer);
  registerStudentTools(mockServer);
  registerRoomTools(mockServer);

  try {
    // 1. Test Login
    console.log('\n--- 1. Testing auth_login ---');
    const loginRes = await tools['auth_login']({ student_id: '22520001', password: 'pass123' });
    console.log(loginRes);

    if (loginRes.isError) {
      console.error('Login failed, cannot continue tests.');
      return;
    }

    const rawData = loginRes.content[0].text;
    const data = JSON.parse(rawData);
    const token = data.data.token;
    console.log('Got Access Token:', token);

    // 2. Test Get Profile
    console.log('\n--- 2. Testing student_get_profile ---');
    const profileRes = await tools['student_get_profile']({ token });
    console.log(profileRes);

    // 3. Test Room Plan Meeting (Composite tool)
    console.log('\n--- 3. Testing room_plan_meeting ---');
    const meetingRes = await tools['room_plan_meeting']({
        token,
        date: "2024-05-20",
        start: "08:00",
        end: "10:00",
        capacity: 50
    });
    console.log(JSON.stringify(meetingRes, null, 2));

    console.log('\nAll tests completed successfully!');

  } catch (error) {
    console.error('Test failed with error:', error);
  }
}

test();
